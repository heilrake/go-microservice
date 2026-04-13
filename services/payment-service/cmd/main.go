package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	paymentInternal "ride-sharing/services/payment-service/internal"
	"ride-sharing/services/payment-service/internal/bootstrap"
	"ride-sharing/services/payment-service/internal/infrastructure/events"
	grpcHandler "ride-sharing/services/payment-service/internal/infrastructure/grpc"
	"ride-sharing/services/payment-service/internal/infrastructure/repository"
	"ride-sharing/services/payment-service/internal/infrastructure/stripe"
	"ride-sharing/services/payment-service/internal/service"
	"ride-sharing/services/payment-service/pkg/types"
	sharedBootstrap "ride-sharing/shared/bootstrap"
	"ride-sharing/shared/env"
	"ride-sharing/shared/tracing"
	sharedTypes "ride-sharing/shared/types"

	grpcserver "google.golang.org/grpc"
)

var grpcAddr = env.GetString("GRPC_ADDR", ":9004")

func main() {
	tracerCfg := tracing.Config{
		ServiceName:    "payment-service",
		Environment:    env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("OTEL_ENDPOINT", "jaeger:4318"),
	}

	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %v", err)
	}

	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	pgConfig := &sharedTypes.PostgresConfig{
		DSN:      env.GetString("DATABASE_URL", "postgres://payment_user:payment_password@payment-postgres:5432/payment_db?sslmode=disable"),
		MaxConns: int32(env.GetInt("DB_MAX_CONNS", 10)),
		MinConns: int32(env.GetInt("DB_MIN_CONNS", 2)),
	}

	if err = sharedBootstrap.RunMigrator(sharedBootstrap.MigratorConfig{
		MigrationsFS:  paymentInternal.Migrations,
		MigrationsDir: "migrations",
		DatabaseURL:   pgConfig.DSN,
		ServiceName:   "payment-service",
	}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	gormDB := sharedBootstrap.InitGorm(pgConfig)
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	appURL := env.GetString("APP_URL", "http://localhost:3000")

	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: env.GetString("STRIPE_SECRET_KEY", ""),
		SuccessURL:      env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success"),
		CancelURL:       env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel"),
	}

	bootstrap.InitStripe(stripeCfg)

	repo := repository.NewPostgresRepository(gormDB)
	paymentProcessor := stripe.NewStripeClient(stripeCfg)
	paymentService := service.NewPaymentService(repo, paymentProcessor)

	// gRPC server
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", grpcAddr, err)
	}

	grpcServer := grpcserver.NewServer(tracing.WithTracingInterceptors()...)
	grpcHandler.NewGRPCHandler(grpcServer, paymentService)

	log.Printf("Starting gRPC server Payment service on %s", grpcAddr)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
			cancel()
		}
	}()

	// RabbitMQ consumers
	rabbitmq := bootstrap.InitRabbitMQ(rabbitMqURI)
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ consumers")

	consumer := events.NewTripConsumer(rabbitmq, paymentService)

	go func() {
		if err := consumer.Listen(); err != nil {
			log.Fatalf("Failed to listen payment queue: %v", err)
		}
	}()

	go func() {
		if err := consumer.ListenCapture(); err != nil {
			log.Fatalf("Failed to listen capture queue: %v", err)
		}
	}()

	go func() {
		if err := consumer.ListenCancel(); err != nil {
			log.Fatalf("Failed to listen cancel queue: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down payment service...")
	grpcServer.GracefulStop()
}
