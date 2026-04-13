package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/shared/auth"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	paymentpb "ride-sharing/shared/proto/payment"
	driverpb "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/user"
	"ride-sharing/shared/tracing"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var tracer = tracing.GetTracer("api-gateway")

// handleTripPreview godoc
// @Summary      Preview trip
// @Description  Calculate route geometry and available fare options before starting a trip
// @Tags         trips
// @Accept       json
// @Produce      json
// @Param        body body previewTripRequest true "Pickup and destination coordinates"
// @Success      201  {object}  contracts.APIResponse  "Route and fare options"
// @Failure      400  {string}  string                 "Invalid JSON"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /trip/preview [post]
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

// handleTripStart godoc
// @Summary      Start trip
// @Description  Create a trip from a previously previewed ride fare
// @Tags         trips
// @Accept       json
// @Produce      json
// @Param        body body startTripRequest true "Selected fare and user ID"
// @Success      201  {object}  contracts.APIResponse  "Created trip ID"
// @Failure      400  {string}  string                 "Invalid JSON"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /trip/start [post]
func handleTripStart(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleTripStart")
		defer span.End()

		var requestBody startTripRequest

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
			return
		}

		trip, err := app.TripClient.Client.CreateTrip(ctx, requestBody.toProto())
		if err != nil {
			log.Printf("Failed to start a trip: %v", err)
			http.Error(w, "Failed to start trip", http.StatusInternalServerError)
			return
		}

		paymentResp, err := app.PaymentClient.Client.CreatePaymentIntent(ctx, &paymentpb.CreatePaymentIntentRequest{
			TripID:        trip.GetTripID(),
			UserID:        trip.GetUserID(),
			AmountInCents: trip.GetAmountInCents(),
			Currency:      trip.GetCurrency(),
		})
		if err != nil {
			log.Printf("Failed to create payment intent for trip %s: %v", trip.TripID, err)
			http.Error(w, "Failed to create payment intent", http.StatusInternalServerError)
			return
		}

		writeJson(w, http.StatusCreated, contracts.APIResponse{
			Data: map[string]string{
				"tripID":       trip.TripID,
				"clientSecret": paymentResp.ClientSecret,
			},
		})
	}
}

// handleTripCancel godoc
// @Summary      Cancel trip
// @Description  Cancel a trip
// @Tags         trips
// @Accept       json
// @Produce      json
// @Param        body body cancelTripRequest true "User ID"
// @Success      200  {object}  contracts.APIResponse  "Cancelled trip"
// @Failure      400  {string}  string                 "Invalid JSON"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /trip/cancel [post]
func handleTripCancel(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleTripCancel")
		defer span.End()

		var requestBody cancelTripRequest

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
			return
		}

		_, err := app.TripClient.Client.CancelTrip(ctx, requestBody.toProto())
		if err != nil {
			log.Printf("Failed to cancel a trip: %v", err)
			http.Error(w, "Failed to cancel trip", http.StatusInternalServerError)
			return
		}

		// Notify the rider via WebSocket that the trip was cancelled
		_ = connManager.SendMessage(requestBody.UserID, contracts.WSMessage{
			Type: contracts.TripEventCancelled,
			Data: map[string]string{"UserID": requestBody.UserID},
		})

		response := contracts.APIResponse{Data: requestBody.UserID}
		fmt.Printf("response %v\n", response)

		writeJson(w, http.StatusOK, response)
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

// handleUserCreation godoc
// @Summary      Create user
// @Description  Register a new user as a rider or driver
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body createUserRequest true "User registration details"
// @Success      201  {object}  contracts.APIResponse  "Created user"
// @Failure      400  {string}  string                 "Invalid JSON or invalid role"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /user/create [post]
func handleUserCreation(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleUserCreation")
		defer span.End()

		var requestBody createUserRequest

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
			return
		}

		role := requestBody.Role
		if role == "" {
			role = "rider"
		}
		if role != "rider" && role != "driver" {
			http.Error(w, "invalid role: must be rider or driver", http.StatusBadRequest)
			return
		}
		user, err := app.UserClient.Client.CreateUser(ctx, &pb.CreateUserRequest{
			Username: requestBody.Username,
			Email:    requestBody.Email,
			Password: requestBody.Password,
			Role:     role,
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

		writeJson(w, http.StatusCreated, response)
	}
}

// handleCreateDriver godoc
// @Summary      Create driver profile
// @Description  Create a driver profile for the authenticated user. Requires driver role JWT.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body createDriverRequest true "Driver name and optional profile picture URL"
// @Success      201  {object}  contracts.APIResponse  "Created driver profile"
// @Failure      400  {string}  string                 "Invalid JSON or missing name"
// @Failure      401  {string}  string                 "Unauthorized"
// @Failure      403  {string}  string                 "Forbidden: driver role required"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /driver [post]
func handleCreateDriver(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleCreateDriver")
		defer span.End()

		userID, err := auth.UserIDFromRequest(r)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		role, _ := auth.RoleFromRequest(r)
		if role != "driver" {
			http.Error(w, "forbidden: driver role required", http.StatusForbidden)
			return
		}

		var req createDriverRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			http.Error(w, "name required", http.StatusBadRequest)
			return
		}

		driver, err := app.DriverClient.Client.CreateDriver(ctx, &driverpb.CreateDriverRequest{
			UserId:         userID,
			Name:           req.Name,
			ProfilePicture: req.ProfilePicture,
		})
		if err != nil {
			log.Printf("Failed to create driver: %v", err)
			http.Error(w, "Failed to create driver", http.StatusInternalServerError)
			return
		}

		writeJson(w, http.StatusCreated, contracts.APIResponse{Data: driver.Driver})
	}
}

// handleGetDriver godoc
// @Summary      Get driver profile
// @Description  Retrieve the driver profile of the authenticated user. Requires driver role JWT.
// @Tags         drivers
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  contracts.APIResponse  "Driver profile"
// @Failure      401  {string}  string                 "Unauthorized"
// @Failure      403  {string}  string                 "Forbidden: driver role required"
// @Failure      404  {string}  string                 "Driver profile not found"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /driver [get]
func handleGetDriver(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleGetDriver")
		defer span.End()

		userID, err := auth.UserIDFromRequest(r)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		role, _ := auth.RoleFromRequest(r)
		if role != "driver" {
			http.Error(w, "forbidden: driver role required", http.StatusForbidden)
			return
		}

		resp, err := app.DriverClient.Client.GetDriver(ctx, &driverpb.GetDriverRequest{UserId: userID})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				http.Error(w, "driver profile not found", http.StatusNotFound)
				return
			}
			log.Printf("Failed to get driver: %v", err)
			http.Error(w, "Failed to get driver", http.StatusInternalServerError)
			return
		}

		writeJson(w, http.StatusOK, contracts.APIResponse{Data: resp.Driver})
	}
}

// handleCreateCar godoc
// @Summary      Add car
// @Description  Register a car for the authenticated driver. Requires driver role JWT.
// @Tags         drivers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body createCarRequest true "Car plate and package slug (e.g. \"economy\", \"comfort\")"
// @Success      201  {object}  contracts.APIResponse  "Created car"
// @Failure      400  {string}  string                 "Missing required fields"
// @Failure      401  {string}  string                 "Unauthorized"
// @Failure      403  {string}  string                 "Forbidden: driver role required"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /driver/cars [post]
func handleCreateCar(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleCreateCar")
		defer span.End()

		userID, err := auth.UserIDFromRequest(r)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		role, _ := auth.RoleFromRequest(r)
		if role != "driver" {
			http.Error(w, "forbidden: driver role required", http.StatusForbidden)
			return
		}

		var req createCarRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to parse JSON", http.StatusBadRequest)
			return
		}
		if req.CarPlate == "" || req.PackageSlug == "" {
			http.Error(w, "car_plate and package_slug required", http.StatusBadRequest)
			return
		}

		car, err := app.DriverClient.Client.CreateCar(ctx, &driverpb.CreateCarRequest{
			UserId:      userID,
			CarPlate:    req.CarPlate,
			PackageSlug: req.PackageSlug,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
				http.Error(w, st.Message(), http.StatusBadRequest)
				return
			}
			log.Printf("Failed to create car: %v", err)
			http.Error(w, "Failed to create car", http.StatusInternalServerError)
			return
		}

		writeJson(w, http.StatusCreated, contracts.APIResponse{Data: car.Car})
	}
}

// handleListCars godoc
// @Summary      List cars
// @Description  Get all cars registered to the authenticated driver. Requires driver role JWT.
// @Tags         drivers
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  contracts.APIResponse  "List of cars"
// @Failure      401  {string}  string                 "Unauthorized"
// @Failure      403  {string}  string                 "Forbidden: driver role required"
// @Failure      500  {string}  string                 "Internal error"
// @Router       /driver/cars [get]
func handleListCars(app *Clients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "handleListCars")
		defer span.End()

		userID, err := auth.UserIDFromRequest(r)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		role, _ := auth.RoleFromRequest(r)
		if role != "driver" {
			http.Error(w, "forbidden: driver role required", http.StatusForbidden)
			return
		}

		resp, err := app.DriverClient.Client.ListCars(ctx, &driverpb.ListCarsRequest{UserId: userID})
		if err != nil {
			// driver profile doesn't exist yet — return empty list
			if st, ok := status.FromError(err); ok && (st.Code() == codes.NotFound || st.Code() == codes.FailedPrecondition) {
				writeJson(w, http.StatusOK, contracts.APIResponse{Data: []*driverpb.Car{}})
				return
			}
			log.Printf("Failed to list cars: %v", err)
			http.Error(w, "Failed to list cars", http.StatusInternalServerError)
			return
		}

		writeJson(w, http.StatusOK, contracts.APIResponse{Data: resp.Cars})
	}
}

// proxyAuth проксує запит на auth-service (login, oauth тощо).
func proxyAuth(authPath string) http.HandlerFunc {
	authBase := env.GetString("AUTH_SERVICE_URL", "http://127.0.0.1:8082")
	targetURL := authBase + authPath
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		proxyReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, targetURL, bytes.NewReader(body))
		if err != nil {
			log.Printf("proxyLogin: NewRequest: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		proxyReq.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			log.Printf("proxyLogin: Do %s: %v", targetURL, err)
			http.Error(w, "auth service unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("proxyLogin: Copy: %v", err)
		}
	}
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
