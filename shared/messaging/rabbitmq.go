package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/tracing"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn      *amqp.Connection
	Channel   *amqp.Channel // used for declaring queues/exchanges and consuming
	publishCh *amqp.Channel // dedicated channel for publishing (separate to avoid race conditions)
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	publishCh, err := conn.Channel()
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to open publish channel: %w", err)
	}

	return &RabbitMQ{
		Conn:      conn,
		Channel:   ch,
		publishCh: publishCh,
	}, nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	log.Printf("Publishing message with routing key: %s", routingKey)

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         jsonMsg,
	}

	return tracing.TracedPublisher(ctx, TripExchange, routingKey, msg, r.publish)
}

func (r *RabbitMQ) publish(ctx context.Context, exchange, routingKey string, msg amqp.Publishing) error {
	return r.publishCh.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		msg)
}

func (r *RabbitMQ) Close() {
	if r.publishCh != nil {
		r.publishCh.Close()
	}
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
