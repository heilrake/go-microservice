package messaging

import (
	"fmt"
	"ride-sharing/shared/contracts"
)

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

	if err := r.DeclareQueue(QueueConfig{
		QueueName:   DriverCmdTripRequestQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.DriverCmdTripRequest},
	}); err != nil {
		return err
	}

	if err := r.DeclareQueue(QueueConfig{
		QueueName:   DriverTripResponseQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline},
	}); err != nil {
		return err
	}

	if err := r.DeclareQueue(QueueConfig{
		QueueName:   NotifyDriverNoDriversFoundQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.TripEventNoDriversFound},
	}); err != nil {
		return err
	}

	if err := r.DeclareQueue(QueueConfig{
		QueueName:   NotifyDriverAssignmentQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.TripEventDriverAssigned},
	}); err != nil {
		return err
	}

	return nil
}
