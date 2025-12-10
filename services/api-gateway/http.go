package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var requestBody previewTripRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	jsonBody, _ := json.Marshal(requestBody)
	reader := bytes.NewBuffer(jsonBody)

	resp, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		http.Error(w, "failed call trip service ", http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	var respBody any
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	fmt.Printf("here\n")

	response := contracts.APIResponse{Data: respBody}

	fmt.Printf("response %v\n", response)

	// surface downstream status code instead of always 200
	writeJson(w, resp.StatusCode, response)
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
