package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	pb "ride-sharing/shared/proto/user"
	"ride-sharing/shared/tracing"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = tracing.GetTracer("api-gateway")

func handleTripPreview(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleTripPreview")
		defer span.End()

		var requestBody previewTripRequest

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
			return
		}

		tripPreview, err := app.TripClient.Client.PreviewTrip(ctx, requestBody.toProto())
		if err != nil {
			log.Printf("Failed to preview a trip: %v", err)
			http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
			return
		}

		response := contracts.APIResponse{Data: tripPreview}

		fmt.Printf("response %v\n", response)

		writeJson(w, http.StatusCreated, response)
	}
}

func handleTripStart(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleTripStart")
		defer span.End()

		var requestBody startTripRequest

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
			return
		}

		tripStart, err := app.TripClient.Client.CreateTrip(ctx, requestBody.toProto())
		if err != nil {
			log.Printf("Failed to start a trip: %v", err)
			http.Error(w, "Failed to start trip", http.StatusInternalServerError)
			return
		}

		response := contracts.APIResponse{Data: tripStart}

		fmt.Printf("response %v\n", response)

		writeJson(w, http.StatusCreated, response)
	}
}

func handleStripeWebhook(w http.ResponseWriter, r *http.Request, rb *messaging.RabbitMQ) {
	ctx, span := tracer.Start(r.Context(), "handleStripeWebhook")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	webhookKey := env.GetString("STRIPE_WEBHOOK_KEY", "")
	if webhookKey == "" {
		log.Printf("Webhook key is required")
		return
	}

	event, err := webhook.ConstructEventWithOptions(
		body,
		r.Header.Get("Stripe-Signature"),
		webhookKey,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	log.Printf("Received Stripe event: %v", event)

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		payload := messaging.PaymentStatusUpdateData{
			TripID:   session.Metadata["trip_id"],
			UserID:   session.Metadata["user_id"],
			DriverID: session.Metadata["driver_id"],
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling payload: %v", err)
			http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
			return
		}

		message := contracts.AmqpMessage{
			OwnerID: session.Metadata["user_id"],
			Data:    payloadBytes,
		}

		if err := rb.PublishMessage(
			ctx,
			contracts.PaymentEventSuccess,
			message,
		); err != nil {
			log.Printf("Error publishing payment event: %v", err)
			http.Error(w, "Failed to publish payment event", http.StatusInternalServerError)
			return
		}
	}
}

func handleUserCreation(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleUserCreation")
		defer span.End()

		var requestBody createUserRequest

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
			return
		}

		user, err := app.UserClient.Client.CreateUser(ctx, &pb.CreateUserRequest{
			Username: requestBody.Username,
			Email:    requestBody.Email,
			Password: requestBody.Password,
		})
		if err != nil {
			if s, ok := status.FromError(err); ok {
				if s.Code() == codes.InvalidArgument {
					http.Error(w, s.Message(), http.StatusBadRequest)
					return
				}
			}
			log.Printf("Failed to create user: %v", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		response := contracts.APIResponse{Data: user}

		fmt.Printf("response %v\n", response)

		writeJson(w, http.StatusCreated, response)
	}
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
