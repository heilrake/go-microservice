package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/services/auth-service/internal"
	"ride-sharing/shared/env"
)

func main() {
	addr := env.GetString("HTTP_ADDR", ":8082")
	log.Println("Starting auth-service on", addr)

	userClient, err := internal.NewUserClient()
	if err != nil {
		log.Fatalf("failed to create user gRPC client: %v", err)
	}
	defer userClient.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/login", internal.EnableCORS(internal.UserLoginHandler(userClient)))
	mux.HandleFunc("POST /driver/login", internal.EnableCORS(internal.DriverLoginHandler(userClient)))
	mux.HandleFunc("POST /auth/oauth", internal.EnableCORS(internal.OAuthHandler(userClient)))
	mux.HandleFunc("OPTIONS /auth/oauth", internal.EnableCORS(func(http.ResponseWriter, *http.Request) {}))

	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("auth-service listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("auth-service shutdown: %v", err)
	}
	log.Println("auth-service stopped")
}
