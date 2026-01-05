package consumers

import (
	"context"
	"log"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentConsumer struct {
	Rabbit *messaging.RabbitMQ
}

func NewPaymentConsumer(r *messaging.RabbitMQ) *PaymentConsumer {
	return &PaymentConsumer{Rabbit: r}
}

func (c *PaymentConsumer) Consume(queue string, handler func(context.Context, amqp.Delivery) error) error {
	if err := c.Rabbit.Channel.Qos(1, 0, false); err != nil {
		return err
	}

	msgs, err := c.Rabbit.Channel.Consume(
		queue,
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := tracing.TracedConsumer(msg, func(ctx context.Context, d amqp.Delivery) error {
				log.Printf("Received message: %s", msg.Body)
				if err := handler(ctx, msg); err != nil {
					log.Printf("Handler error: %v", err)
					_ = msg.Nack(false, false)
					return err
				}
				_ = msg.Ack(false)
				return nil
			}); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}()

	return nil
}
