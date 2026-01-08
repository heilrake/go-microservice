package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	driverInternal "ride-sharing/services/driver-service/internal"
	"ride-sharing/services/driver-service/internal/events"
	infrastructure "ride-sharing/services/driver-service/internal/infrastructure/db"
	grpcHandler "ride-sharing/services/driver-service/internal/infrastructure/grpc"
	"ride-sharing/services/driver-service/internal/infrastructure/repository"
	"ride-sharing/services/driver-service/internal/service"
	sharedBootstrap "ride-sharing/shared/bootstrap"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"
	"ride-sharing/shared/types"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr string = ":9094"

func main() {
	// Initialize Tracing
	tracerCfg := tracing.Config{
		ServiceName:    "driver-service",
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
		DSN:      env.GetString("DATABASE_URL", "postgres://driver_user:driver_password@driver-postgres:5432/driver_db?sslmode=disable"),
		MaxConns: int32(env.GetInt("DB_MAX_CONNS", 10)),
		MinConns: int32(env.GetInt("DB_MIN_CONNS", 2)),
	}

	// Run migrations on startup
	err = sharedBootstrap.RunMigrator(sharedBootstrap.MigratorConfig{
		MigrationsFS:  driverInternal.Migrations,
		MigrationsDir: "migrations",
		DatabaseURL:   pgConfig.DSN,
		ServiceName:   "driver-service",
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
	svc := service.NewDriverService(repo)

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
	grpcServer := grpcserver.NewServer(tracing.WithTracingInterceptors()...)

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	consumer := events.NewTripConsumer(rabbitmq, svc)
	go func() {
		if err := consumer.Listen(); err != nil {
			log.Fatalf("Failed to listen to the message: %v", err)
		}
	}()

	// Create the service and register the gRPC handler

	grpcHandler.NewGrpcHandler(grpcServer, svc)

	log.Printf("Starting gRPC server Driver service on port %s", lis.Addr().String())

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
