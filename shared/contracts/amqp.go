package contracts

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string `json:"ownerId"`
	Data    []byte `json:"data"`
}

// Routing keys - using consistent event/command patterns
const (
	// Trip events (trip.event.*)
	TripEventCreated             = "trip.event.created"
	TripEventCancelled           = "trip.event.cancelled"
	TripEventCompleted           = "trip.event.completed"
	TripEventDriverAssigned      = "trip.event.driver_assigned"
	TripEventNoDriversFound      = "trip.event.no_drivers_found"
	TripEventDriverNotInterested = "trip.event.driver_not_interested"

	// Rider commands (rider.cmd.*)
	RiderCmdPaymentConfirmed = "rider.cmd.payment_confirmed"

	// Driver commands (driver.cmd.*)
	DriverCmdTripComplete = "driver.cmd.trip_complete"
	DriverCmdTripRequest  = "driver.cmd.trip_request"
	DriverCmdTripAccept   = "driver.cmd.trip_accept"
	DriverCmdTripDecline  = "driver.cmd.trip_decline"
	DriverCmdLocation     = "driver.cmd.location"
	DriverCmdRegister     = "driver.cmd.register"

	// Driver events (driver.event.*)
	DriverEventDriverNotified = "driver.event.driver_notified"

	DriverCmdTripRequestExpired = "driver.cmd.trip_request_expired"

	// Payment events (payment.event.*)
	PaymentEventSessionCreated = "payment.event.session_created"
	PaymentEventSuccess        = "payment.event.success"
	PaymentEventFailed         = "payment.event.failed"
	PaymentEventCancelled      = "payment.event.cancelled"

	// Payment commands (payment.cmd.*)
	PaymentCmdCreateSession  = "payment.cmd.create_session"
	PaymentCmdCapturePayment = "payment.cmd.capture_payment"
	PaymentCmdCancelPayment  = "payment.cmd.cancel_payment"

	// User events (user.event.*)
	UserEventCreated = "user.event.created"
	UserEventUpdated = "user.event.updated"
	UserEventDeleted = "user.event.deleted"
)
