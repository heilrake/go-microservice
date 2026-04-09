package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/messaging/consumers"
	pbd "ride-sharing/shared/proto/driver"

	"github.com/rabbitmq/amqp091-go"
)

type DriverConsumer struct {
	rabbitmq    *messaging.RabbitMQ
	service     domain.TripService
	retryCounts sync.Map // tripID → int
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, service domain.TripService) *DriverConsumer {
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
		default:
			log.Printf("unknown trip event: %+v", payload)
		}

		return nil
	})
}

// ListenForDriverNotified listens for the event published by driver-service when a driver
// has been selected and notified, and updates the trip status to "awaiting_driver".
func (c *DriverConsumer) ListenForDriverNotified() error {
	return consumers.NewFindDriversConsumer(c.rabbitmq).Consume(messaging.DriverNotifiedQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal driver_notified message: %v", err)
			return err
		}

		var payload messaging.DriverNotifiedData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal driver_notified payload: %v", err)
			return err
		}

		log.Printf("Driver notified for trip %s — updating status to awaiting_driver", payload.TripID)

		if err := c.service.UpdateTrip(ctx, payload.TripID, "awaiting_driver", nil); err != nil {
			log.Printf("Failed to update trip %s status: %v", payload.TripID, err)
			return err
		}

		return nil
	})
}

// ListenForNoDriversFound listens for the event published by driver-service when
// no drivers could be found, and marks the trip as "cancelled" in the DB.
func (c *DriverConsumer) ListenForNoDriversFound() error {
	return consumers.NewFindDriversConsumer(c.rabbitmq).Consume(messaging.TripSearchFailedQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal no_drivers_found message: %v", err)
			return err
		}

		var payload messaging.NoDriversFoundData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal no_drivers_found payload: %v", err)
			return err
		}

		log.Printf("No drivers found for trip %s — updating status to cancelled", payload.TripID)

		c.retryCounts.Delete(payload.TripID)

		if err := c.service.UpdateTrip(ctx, payload.TripID, "cancelled", nil); err != nil {
			log.Printf("Failed to update trip %s status to cancelled: %v", payload.TripID, err)
			return err
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

	// Guard: if trip is already assigned (e.g. timeout fired after driver accepted), skip
	if trip.Status == "assigned" {
		log.Printf("Trip %s is already assigned — ignoring duplicate accept", tripID)
		return nil
	}

	c.retryCounts.Delete(tripID)

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

	marshalledPayload, err := json.Marshal(messaging.PaymentTripResponseData{
		TripID:   tripID,
		UserID:   trip.UserID,
		DriverID: driver.Id,
		Amount:   int64(trip.RideFare.TotalPriceInCents),
		Currency: "USD",
	})
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.PaymentCmdCreateSession,
		contracts.AmqpMessage{
			OwnerID: trip.UserID,
			Data:    marshalledPayload,
		},
	); err != nil {
		return err
	}

	return nil
}

func (c *DriverConsumer) handleTripDeclined(ctx context.Context, tripID string, riderID string) error {
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	// Guard: if trip is already assigned (race between decline and accept), skip re-search
	if trip.Status == "assigned" {
		log.Printf("Trip %s is already assigned — ignoring late decline", tripID)
		return nil
	}

	retryCount := 0
	if v, ok := c.retryCounts.Load(tripID); ok {
		retryCount = v.(int)
	}
	retryCount++
	c.retryCounts.Store(tripID, retryCount)

	log.Printf("Trip %s declined — retry attempt %d", tripID, retryCount)

	newPayload := messaging.TripEventData{
		Trip:       trip.ToProto(),
		RetryCount: retryCount,
	}

	marshalledEvent, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverNotInterested, contracts.AmqpMessage{
		OwnerID: riderID,
		Data:    marshalledEvent,
	}); err != nil {
		return err
	}

	return nil
}
