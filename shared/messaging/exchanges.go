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

	// Queue for driver-service to find available drivers when trip is created
	if err := r.DeclareQueue(QueueConfig{
		QueueName:   FindAvailableDriversQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.TripEventCreated, contracts.TripEventDriverNotInterested},
	}); err != nil {
		return err
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

	// Driver-service listens here to cancel pending timers when driver responds
	if err := r.DeclareQueue(QueueConfig{
		QueueName:   DriverTripAckQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline},
	}); err != nil {
		return err
	}

	// Trip-service listens here to update trip status to "awaiting_driver"
	if err := r.DeclareQueue(QueueConfig{
		QueueName:   DriverNotifiedQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.DriverEventDriverNotified},
	}); err != nil {
		return err
	}

	// API Gateway routes this to the driver via WebSocket when their request times out
	if err := r.DeclareQueue(QueueConfig{
		QueueName:   DriverTripRequestExpiredQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.DriverCmdTripRequestExpired},
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

	// Trip-service listens here to update trip status when no drivers are found
	if err := r.DeclareQueue(QueueConfig{
		QueueName:   TripSearchFailedQueue,
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

	if err := r.DeclareQueue(QueueConfig{
		QueueName:   NotifyPaymentSessionCreatedQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.PaymentEventSessionCreated},
	}); err != nil {
		return err
	}

	if err := r.DeclareQueue(QueueConfig{
		QueueName:   PaymentTripResponseQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.PaymentCmdCreateSession},
	}); err != nil {
		return err
	}

	if err := r.DeclareQueue(QueueConfig{
		QueueName:   NotifyPaymentSuccessQueue,
		Exchanges:   []string{TripExchange},
		RoutingKeys: []string{contracts.PaymentEventSuccess},
	}); err != nil {
		return err
	}

	return nil
}
