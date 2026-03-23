package publishers

import (
	"context"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
)

type TripPublisher struct {
	Rabbit *messaging.RabbitMQ
}

func NewTripPublisher(r *messaging.RabbitMQ) *TripPublisher {
	return &TripPublisher{Rabbit: r}
}

func (p *TripPublisher) Publish(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	return p.Rabbit.PublishMessage(ctx, routingKey, message)
}
