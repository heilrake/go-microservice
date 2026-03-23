package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/proto/driver"
)

var (
	connManager = messaging.NewConnectionManager()
)

func handleRiderWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := connManager.Upgrade(w, r)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("No user ID provided")
		return
	}

	connManager.Add(userID, conn)
	defer connManager.Remove(userID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", message)
	}
}

func handleDriverWebSocket(w http.ResponseWriter, r *http.Request, rabbitmq *messaging.RabbitMQ) {
	conn, err := connManager.Upgrade(w, r)

	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("No user ID provided")
		return
	}

	carID := r.URL.Query().Get("carID")
	if carID == "" {
		log.Println("No carID provided")
		return
	}

	lat, _ := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)

	// Add connection to manager
	connManager.Add(userID, conn)

	ctx := r.Context()

	driverService, err := grpc_client.NewDriverServiceClient()
	if err != nil {
		log.Printf("Failed to create driver gRPC client: %v", err)
		_ = connManager.SendMessage(userID, contracts.WSMessage{Type: "driver.cmd.error", Data: "service unavailable"})
		return
	}
	defer func() {
		connManager.Remove(userID)
		driverService.Client.UnRegisterDriver(ctx, &driver.RegisterDriverRequest{DriverID: userID})
		driverService.Close()
		log.Println("Driver unregistered: ", userID)
	}()

	driverData, err := driverService.Client.RegisterDriver(ctx, &driver.RegisterDriverRequest{
		DriverID:  userID,
		CarId:     carID,
		Latitude:  lat,
		Longitude: lon,
	})
	if err != nil {
		log.Printf("Error registering driver (userID=%s carID=%s): %v", userID, carID, err)
		_ = connManager.SendMessage(userID, contracts.WSMessage{Type: "driver.cmd.error", Data: err.Error()})
		return
	}

	if err := connManager.SendMessage(userID, contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driverData.Driver,
	}); err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		type driverMessage struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}

		var driverMsg driverMessage
		if err := json.Unmarshal(message, &driverMsg); err != nil {
			log.Printf("Error unmarshaling driver message: %v", err)
			continue
		}

		// Handle the different message type
		switch driverMsg.Type {
		case contracts.DriverCmdLocation:
			// Handle driver location update in the future
			continue
		case contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline:
			// Forward the message to RabbitMQ
			if err := rabbitmq.PublishMessage(
				ctx,
				driverMsg.Type,
				contracts.AmqpMessage{
					OwnerID: userID,
					Data:    driverMsg.Data,
				},
			); err != nil {
				log.Printf("Error publishing message to RabbitMQ: %v", err)
			}
		default:
			log.Printf("Unknown message type: %s", driverMsg.Type)
		}
	}
}
