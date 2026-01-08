package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/driver-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/messaging/consumers"
	tripPublishers "ride-sharing/shared/messaging/publishers"

	"math/rand"

	"github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  domain.DriverService
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ, service domain.DriverService) *tripConsumer {
	return &tripConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *tripConsumer) Listen() error {
	return consumers.NewFindDriversConsumer(c.rabbitmq).Consume(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		log.Printf("driver received message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
			return c.handleFindAndNotifyDrivers(ctx, payload)
		}

		log.Printf("unknown trip event: %+v", payload)

		return nil
	})
}

func (c *tripConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	if payload.Trip == nil || payload.Trip.SelectedFare == nil {
		log.Printf("Trip has no selected fare, cannot find drivers")
		return fmt.Errorf("trip has no selected fare")
	}

	suitableDrivers, err := c.service.FindAvailableDrivers(ctx, payload.Trip.SelectedFare.PackageSlug)
	if err != nil {
		return err
	}

	log.Printf("Found suitable drivers %v", len(suitableDrivers))

	if len(suitableDrivers) == 0 {
		// Notify the driver that no drivers are available
		if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
		}); err != nil {
			log.Printf("Failed to publish message to exchange: %v", err)
			return err
		}

		return nil
	}

	// Get a random index from the matching drivers
	randomIndex := rand.Intn(len(suitableDrivers))

	suitableDriver := suitableDrivers[randomIndex]

	marshalledEvent, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Notify the driver about a potential trip
	if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: suitableDriver.ID,
		Data:    marshalledEvent,
	}); err != nil {
		log.Printf("Failed to publish message to exchange: %v", err)
		return err
	}

	return nil
}
