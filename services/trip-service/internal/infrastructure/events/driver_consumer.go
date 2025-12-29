package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/messaging/consumers"
	pbd "ride-sharing/shared/proto/driver"

	"github.com/rabbitmq/amqp091-go"
)

type DriverConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  service.TripService
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, service service.TripService) *DriverConsumer {
	return &DriverConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *DriverConsumer) Listen() error {
	return consumers.NewFindDriversConsumer(c.rabbitmq).Consume(messaging.DriverTripResponseQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		log.Printf("driver response received message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			if err := c.handleTripAccepted(ctx, payload.TripID, payload.Driver); err != nil {
				log.Printf("Failed to handle the trip accept: %v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			if err := c.handleTripDeclined(ctx, payload.TripID, payload.RiderID); err != nil {
				log.Printf("Failed to handle the trip decline: %v", err)
				return err
			}
			return nil
		default:
			log.Printf("unknown trip event: %+v", payload)
		}

		return nil
	})
}

func (c *DriverConsumer) handleTripAccepted(ctx context.Context, tripID string, driver *pbd.Driver) error {
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	if trip == nil {
		return fmt.Errorf("Trip was not found %s", tripID)
	}

	if err := c.service.UpdateTrip(ctx, tripID, "assigned", driver); err != nil {
		return err
	}

	trip, err = c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	marshalledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	// Notify the rider that a driver has been assigned
	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    marshalledTrip,
	}); err != nil {
		return err
	}

	return nil
}

func (c *DriverConsumer) handleTripDeclined(ctx context.Context, tripID string, riderID string) error {
	// When a driver declines, we should try to find another driver

	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	newPayload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	marshalledEvent, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	// Notify the rider that the trip has been declined
	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverNotInterested, contracts.AmqpMessage{
		OwnerID: riderID,
		Data:    marshalledEvent,
	}); err != nil {
		return err
	}

	return nil
}
