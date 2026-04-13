package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"

	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/messaging/consumers"
)

type TripConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  domain.Service
}

func NewTripConsumer(rabbitmq *messaging.RabbitMQ, service domain.Service) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *TripConsumer) Listen() error {
	return consumers.NewTripPaymentConsumer(c.rabbitmq).Consume(messaging.PaymentTripResponseQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		switch msg.RoutingKey {
		case contracts.PaymentCmdCreateSession:
			var payload messaging.PaymentTripResponseData
			if err := json.Unmarshal(message.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal payload: %v", err)
				return err
			}
			if err := c.handleCreatePaymentIntent(ctx, payload); err != nil {
				log.Printf("Failed to handle create payment intent: %v", err)
				return err
			}
		}

		return nil
	})
}

func (c *TripConsumer) ListenCapture() error {
	return consumers.NewTripPaymentConsumer(c.rabbitmq).Consume(messaging.PaymentCaptureQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal capture message: %v", err)
			return err
		}

		var payload messaging.TripCompletedData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal capture payload: %v", err)
			return err
		}

		log.Printf("Capturing payment for trip: %s", payload.TripID)
		if err := c.service.CapturePayment(ctx, payload.TripID); err != nil {
			log.Printf("Failed to capture payment for trip %s: %v", payload.TripID, err)
			return err
		}

		log.Printf("Payment captured for trip: %s", payload.TripID)
		return nil
	})
}

func (c *TripConsumer) ListenCancel() error {
	return consumers.NewTripPaymentConsumer(c.rabbitmq).Consume(messaging.PaymentCancelQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal cancel message: %v", err)
			return err
		}

		var payload messaging.NoDriversFoundData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal cancel payload: %v", err)
			return err
		}

		log.Printf("Cancelling payment for trip: %s", payload.TripID)
		if err := c.service.CancelPayment(ctx, payload.TripID); err != nil {
			log.Printf("Failed to cancel payment for trip %s: %v", payload.TripID, err)
			return err
		}

		log.Printf("Payment cancelled for trip: %s", payload.TripID)
		return nil
	})
}

func (c *TripConsumer) handleCreatePaymentIntent(ctx context.Context, payload messaging.PaymentTripResponseData) error {
	log.Printf("Creating payment intent for trip: %s", payload.TripID)

	intent, err := c.service.CreatePaymentIntent(
		ctx,
		payload.TripID,
		payload.UserID,
		int64(payload.Amount),
		payload.Currency,
	)
	if err != nil {
		log.Printf("Failed to create payment intent: %v", err)
		return err
	}

	log.Printf("Payment intent created: %s for trip: %s", intent.StripePaymentIntentID, intent.TripID)
	return nil
}
