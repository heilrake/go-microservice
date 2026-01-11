package publishers

import (
	"context"
	"encoding/json"
	"fmt"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UserPublisher struct {
	Rabbit *messaging.RabbitMQ
}

func NewUserPublisher(r *messaging.RabbitMQ) *UserPublisher {
	return &UserPublisher{Rabbit: r}
}

func (p *UserPublisher) Publish(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.Rabbit.Channel.PublishWithContext(
		ctx,
		messaging.UserExchange,
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

