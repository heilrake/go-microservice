package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"
)

var (
	httpAddr    = env.GetString("HTTP_ADDR", ":8081")
	rabbitMqURI = env.GetString("RABBITMQ_URI", "amqp://admin:admin@rabbitmq:5672")
)

func main() {
	log.Println("Starting API Gateway")

	// Initialize Tracing
	tracerCfg := tracing.Config{
		ServiceName:    "api-gateway",
		Environment:    env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("OTEL_ENDPOINT", "jaeger:4318"),
	}

	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	// Initialize shared gRPC clients (один раз при старті!)
	app, err := InitializeClients()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer app.CloseAllClients()

	log.Println("Initialized gRPC clients")

	mux := http.NewServeMux()

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	// HTTP routes - передаємо app для доступу до shared clients
	mux.Handle("POST /trip/preview", tracing.WrapHandlerFunc(enableCORS(handleTripPreview(app)), "handleTripPreview"))
	mux.Handle("POST /trip/start", tracing.WrapHandlerFunc(enableCORS(handleTripStart(app)), "handleTripStart"))
	mux.Handle("POST /user/create", tracing.WrapHandlerFunc(enableCORS(handleUserCreation(app)), "handleUserCreation"))
	mux.Handle("/ws/drivers", tracing.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleDriverWebSocket(w, r, rabbitmq)
	}, "/ws/drivers"))
	mux.Handle("/ws/riders", tracing.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleRiderWebSocket(w, r, rabbitmq)
	}, "/ws/riders"))

	mux.Handle("/webhook/stripe", tracing.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleStripeWebhook(w, r, rabbitmq)
	}, "/webhook/stripe"))

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Serve listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)
	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not  stop the server graceful: %v", err)
			server.Close()
		}
	}

}
