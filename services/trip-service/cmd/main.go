package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	tripinternal "ride-sharing/services/trip-service/internal"
	"ride-sharing/services/trip-service/internal/infrastructure/db"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	service "ride-sharing/services/trip-service/internal/service"
	sharedBootstrap "ride-sharing/shared/bootstrap"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"
	"ride-sharing/shared/types"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr string = ":9093"

func main() {
	// Initialize Tracing
	tracerCfg := tracing.Config{
		ServiceName:    "trip-service",
		Environment:    env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("OTEL_ENDPOINT", "jaeger:4318"),
	}
	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %v", err)
	}

	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://admin:admin@rabbitmq:5672")

	// Initialize PostgreSQL config
	pgConfig := &types.PostgresConfig{
		DSN:      env.GetString("DATABASE_URL", "postgres://trip_user:trip_password@trip-postgres:5432/trip_db?sslmode=disable"),
		MaxConns: int32(env.GetInt("DB_MAX_CONNS", 10)),
		MinConns: int32(env.GetInt("DB_MIN_CONNS", 2)),
	}

	// Run migrations on startup
	err = sharedBootstrap.RunMigrator(sharedBootstrap.MigratorConfig{
		MigrationsFS:  tripinternal.Migrations,
		MigrationsDir: "migrations",
		DatabaseURL:   pgConfig.DSN,
		ServiceName:   "trip-service",
	})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Initialize GORM database connection
	gormDB := infrastructure.InitGorm(pgConfig)
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	// Initialize PostgreSQL repository with GORM
	repo := repository.NewPostgresRepository(gormDB)
	svc := service.NewTripServer(repo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	if err := rabbitmq.DeclareExchanges(); err != nil {
		log.Fatalf("Failed to declare exchanges: %v", err)
	}

	publisher := events.NewTripEventPublisher(rabbitmq)

	driverConsumer := events.NewDriverConsumer(rabbitmq, svc)
	go func() {
		if err := driverConsumer.Listen(); err != nil {
			log.Fatalf("Failed to listen to driver messages: %v", err)
		}
	}()

	grpcServer := grpcserver.NewServer(tracing.WithTracingInterceptors()...)
	grpc.NewGRPCHandler(grpcServer, svc, publisher)

	log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
