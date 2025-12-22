package publishers

import (
	"context"
	"encoding/json"
	"fmt"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TripPublisher struct {
	Rabbit *messaging.RabbitMQ
}

func NewTripPublisher(r *messaging.RabbitMQ) *TripPublisher {
	return &TripPublisher{Rabbit: r}
}

func (p *TripPublisher) Publish(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.Rabbit.Channel.PublishWithContext(
		ctx,
		messaging.TripExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
