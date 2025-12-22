package messaging

import "fmt"

const (
	TripExchange = "trip"
	UserExchange = "user"
)

func (r *RabbitMQ) DeclareExchanges() error {
	exchanges := []struct {
		Name string
		Type string
	}{
		{TripExchange, "topic"},
		{UserExchange, "topic"},
	}

	for _, ex := range exchanges {
		if err := r.Channel.ExchangeDeclare(
			ex.Name,
			ex.Type,
			true,  // durable
			false, // auto-delete
			false, // internal
			false, // no-wait
			nil,   // args
		); err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", ex.Name, err)
		}
	}

	return nil
}
