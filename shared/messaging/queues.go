package messaging

import "fmt"

func (r *RabbitMQ) DeclareQueue(queueName string, exchange string, routingKeys []string) error {
	q, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare failed: %w", err)
	}

	for _, key := range routingKeys {
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

	return nil
}
