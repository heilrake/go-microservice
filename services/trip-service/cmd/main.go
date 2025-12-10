package main

import (
	"log"
	"net/http"
	internalHttp "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
)

func main() {

	inmem := repository.NewInmemRepository()
	service := service.NewTripServer(inmem)

	mux := http.NewServeMux()

	internalHttpHandler := internalHttp.HttpHandler{Service: service}

	mux.HandleFunc("POST /preview", internalHttpHandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("API Gateway have error: %v", err)
	}

}
