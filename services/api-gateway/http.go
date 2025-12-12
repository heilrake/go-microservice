package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	grpc_client "ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var requestBody previewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	tripService, err := grpc_client.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(r.Context(), requestBody.toProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
	}

	response := contracts.APIResponse{Data: tripPreview}

	fmt.Printf("response %v\n", response)

	writeJson(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var requestBody startTripRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	tripService, err := grpc_client.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	tripStart, err := tripService.Client.CreateTrip(r.Context(), requestBody.toProto())
	if err != nil {
		log.Printf("Failed to start a trip: %v", err)
		http.Error(w, "Failed to start trip", http.StatusInternalServerError)
	}

	response := contracts.APIResponse{Data: tripStart}

	fmt.Printf("response %v\n", response)

	writeJson(w, http.StatusCreated, response)
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
