package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"ride-sharing/services/driver-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/messaging/consumers"
	tripPublishers "ride-sharing/shared/messaging/publishers"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type tripConsumer struct {
	rabbitmq        *messaging.RabbitMQ
	service         domain.DriverService
	pendingRequests sync.Map // tripID → context.CancelFunc
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

// ListenForAck listens on DriverTripAckQueue to cancel pending timers when a driver responds.
func (c *tripConsumer) ListenForAck() error {
	return consumers.NewFindDriversConsumer(c.rabbitmq).Consume(messaging.DriverTripAckQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("[ack] Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("[ack] Failed to unmarshal payload: %v", err)
			return err
		}

		if cancel, ok := c.pendingRequests.LoadAndDelete(payload.TripID); ok {
			cancel.(context.CancelFunc)()
			log.Printf("[ack] Timer cancelled for trip %s — driver responded", payload.TripID)
		}

		return nil
	})
}

const maxSearchAttempts = 1

func (c *tripConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	tracer := otel.GetTracerProvider().Tracer("driver-service")
	ctx, span := tracer.Start(ctx, "driver.find_and_notify",
		trace.WithAttributes(
			attribute.String("trip.id", payload.Trip.GetId()),
			attribute.String("trip.user_id", payload.Trip.GetUserID()),
			attribute.String("trip.package_slug", payload.Trip.GetSelectedFare().GetPackageSlug()),
			attribute.Int("trip.retry_count", payload.RetryCount),
		),
	)
	defer span.End()

	if payload.RetryCount >= maxSearchAttempts {
		log.Printf("Trip %s exceeded max search attempts (%d) — giving up", payload.Trip.GetId(), maxSearchAttempts)
		span.AddEvent("max_attempts_exceeded")
		noDriversData, err := json.Marshal(messaging.NoDriversFoundData{
			TripID:  payload.Trip.GetId(),
			RiderID: payload.Trip.GetUserID(),
		})
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
			Data:    noDriversData,
		}); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		return nil
	}

	if payload.Trip == nil || payload.Trip.SelectedFare == nil {
		span.SetStatus(codes.Error, "trip has no selected fare")
		log.Printf("Trip has no selected fare, cannot find drivers")
		return fmt.Errorf("trip has no selected fare")
	}

	// --- Search ---
	_, searchSpan := tracer.Start(ctx, "driver.search",
		trace.WithAttributes(attribute.String("trip.package_slug", payload.Trip.SelectedFare.PackageSlug)),
	)
	suitableDrivers, err := c.service.FindAvailableDrivers(ctx, payload.Trip.SelectedFare.PackageSlug)
	if err != nil {
		searchSpan.SetStatus(codes.Error, err.Error())
		searchSpan.End()
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	searchSpan.SetAttributes(attribute.Int("drivers.found", len(suitableDrivers)))
	searchSpan.End()

	log.Printf("Found suitable drivers %v", len(suitableDrivers))

	if len(suitableDrivers) == 0 {
		span.AddEvent("no_drivers_found")
		noDriversData, err := json.Marshal(messaging.NoDriversFoundData{
			TripID:  payload.Trip.GetId(),
			RiderID: payload.Trip.GetUserID(),
		})
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
			Data:    noDriversData,
		}); err != nil {
			log.Printf("Failed to publish message to exchange: %v", err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}

		return nil
	}

	randomIndex := rand.Intn(len(suitableDrivers))
	suitableDriver := suitableDrivers[randomIndex]

	if suitableDriver.UserID == "" {
		log.Printf("Selected driver %s has empty user_id — skipping (stale DB row, recreate driver profile)", suitableDriver.ID)
		span.SetStatus(codes.Error, "selected driver has no user_id")
		span.SetAttributes(attribute.String("driver.id", suitableDriver.ID))
		return fmt.Errorf("selected driver has no user_id")
	}

	span.SetAttributes(
		attribute.String("driver.id", suitableDriver.ID),
		attribute.String("driver.user_id", suitableDriver.UserID),
	)

	marshalledEvent, err := json.Marshal(payload)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// --- Notify driver ---
	_, notifySpan := tracer.Start(ctx, "driver.notify",
		trace.WithAttributes(
			attribute.String("driver.user_id", suitableDriver.UserID),
			attribute.String("trip.id", payload.Trip.Id),
		),
	)
	if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: suitableDriver.UserID,
		Data:    marshalledEvent,
	}); err != nil {
		log.Printf("Failed to publish message to exchange: %v", err)
		notifySpan.SetStatus(codes.Error, err.Error())
		notifySpan.End()
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	notifySpan.End()

	// Notify trip-service to update trip status to "awaiting_driver"
	notifiedData, err := json.Marshal(messaging.DriverNotifiedData{
		TripID:  payload.Trip.Id,
		RiderID: payload.Trip.UserID,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(ctx, contracts.DriverEventDriverNotified, contracts.AmqpMessage{
		OwnerID: payload.Trip.UserID,
		Data:    notifiedData,
	}); err != nil {
		log.Printf("Failed to publish driver_notified event: %v", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Cancel any existing timer for this trip (handles re-try after decline)
	if existingCancel, loaded := c.pendingRequests.LoadAndDelete(payload.Trip.Id); loaded {
		existingCancel.(context.CancelFunc)()
	}

	// Capture span context to carry the trace into the goroutine
	spanCtx := trace.SpanContextFromContext(ctx)

	// Start 15-second timeout: if driver doesn't respond, try next driver
	timeoutCtx, cancel := context.WithCancel(context.Background())
	c.pendingRequests.Store(payload.Trip.Id, cancel)

	go func(tripID string, p messaging.TripEventData, driverUserID string) {
		select {
		case <-time.After(15 * time.Second):
			c.pendingRequests.Delete(tripID)
			log.Printf("Driver timeout for trip %s (attempt %d) — notifying driver and searching for next", tripID, p.RetryCount+1)

			// Link the timeout span back to the original request span
			linkedCtx := trace.ContextWithRemoteSpanContext(context.Background(), spanCtx)
			_, timeoutSpan := tracer.Start(linkedCtx, "driver.response_timeout",
				trace.WithLinks(trace.Link{SpanContext: spanCtx}),
				trace.WithAttributes(
					attribute.String("trip.id", tripID),
					attribute.String("driver.user_id", driverUserID),
					attribute.Float64("timeout.seconds", 15),
					attribute.Int("trip.retry_count", p.RetryCount),
				),
			)

			if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(
				linkedCtx,
				contracts.DriverCmdTripRequestExpired,
				contracts.AmqpMessage{
					OwnerID: driverUserID,
				},
			); err != nil {
				log.Printf("Failed to publish trip_request_expired to driver %s: %v", driverUserID, err)
				timeoutSpan.SetStatus(codes.Error, err.Error())
			}

			p.RetryCount++
			marshalledTimeout, err := json.Marshal(p)
			if err != nil {
				log.Printf("Failed to marshal timeout payload: %v", err)
				timeoutSpan.SetStatus(codes.Error, err.Error())
				timeoutSpan.End()
				return
			}

			if err := tripPublishers.NewTripPublisher(c.rabbitmq).Publish(
				linkedCtx,
				contracts.TripEventDriverNotInterested,
				contracts.AmqpMessage{
					OwnerID: p.Trip.UserID,
					Data:    marshalledTimeout,
				},
			); err != nil {
				log.Printf("Failed to publish timeout event for trip %s: %v", tripID, err)
				timeoutSpan.SetStatus(codes.Error, err.Error())
			}

			timeoutSpan.End()
		case <-timeoutCtx.Done():
			// Driver responded (accept or decline) — timer was cancelled via ListenForAck
		}
	}(payload.Trip.Id, payload, suitableDriver.UserID)

	return nil
}
