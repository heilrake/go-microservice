package consumers

import (
	"context"
	"log"
	"ride-sharing/shared/messaging"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TripPaymentConsumer struct {
	Rabbit *messaging.RabbitMQ
}

func NewTripPaymentConsumer(r *messaging.RabbitMQ) *TripPaymentConsumer {
	return &TripPaymentConsumer{Rabbit: r}
}

func (c *TripPaymentConsumer) Consume(queue string, handler func(context.Context, amqp.Delivery) error) error {
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

	ctx := context.Background()
	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.Body)
			if err := handler(ctx, msg); err != nil {
				log.Printf("Handler error: %v", err)
				_ = msg.Nack(false, false)
				continue
			}
			_ = msg.Ack(false)
		}
	}()

	return nil
}
