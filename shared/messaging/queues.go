package messaging

import "fmt"

type QueueConfig struct {
	QueueName   string
	Exchanges   []string
	RoutingKeys []string
}

func (r *RabbitMQ) DeclareQueue(config QueueConfig) error {
	q, err := r.Channel.QueueDeclare(
		config.QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
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
