package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type QueueConfig struct {
	QueueName   string
	Exchanges   []string
	RoutingKeys []string
	//DLQ
	DeadLetterQueue string
	MaxRetries      int
}

func (r *RabbitMQ) DeclareQueue(config QueueConfig) error {
	args := amqp091.Table{}
	if config.DeadLetterQueue != "" {
		args["x-queue-type"] = "quorum"
		args["x-dead-letter-exchange"] = ""
		args["x-dead-letter-routing-key"] = config.DeadLetterQueue
		if config.MaxRetries > 0 {
			args["x-delivery-limit"] = int64(config.MaxRetries)
		}
	}
	q, err := r.Channel.QueueDeclare(
		config.QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	for _, exchange := range config.Exchanges {
		for _, key := range config.RoutingKeys {
			if err := r.Channel.QueueBind(
				q.Name,
				key,
				exchange,
				false,
				nil,
			); err != nil {
				return fmt.Errorf("queue bind failed: %w", err)
			}
		}
	}

	return nil
}
