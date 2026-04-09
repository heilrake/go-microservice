package messaging

import (
	pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
	pbu "ride-sharing/shared/proto/user"
)

const (
	FindAvailableDriversQueue        = "find_available_drivers"
	DriverCmdTripRequestQueue        = "driver_cmd_trip_request"
	DriverTripResponseQueue          = "driver_trip_response"
	DriverTripAckQueue               = "driver_trip_ack"
	DriverNotifiedQueue              = "driver_notified"
	DriverTripRequestExpiredQueue    = "driver_trip_request_expired"
	NotifyDriverNoDriversFoundQueue  = "notify_driver_no_drivers_found"
	TripSearchFailedQueue            = "trip_search_failed"
	NotifyDriverAssignmentQueue      = "notify_driver_assignment"
	PaymentTripResponseQueue         = "payment_trip_response"
	NotifyPaymentSessionCreatedQueue = "notify_payment_session_created"
	NotifyPaymentSuccessQueue        = "notify_payment_success"
)

type TripEventData struct {
	Trip       *pb.Trip `json:"trip"`
	RetryCount int      `json:"retryCount"`
}

type DriverTripResponseData struct {
	Driver  *pbd.Driver `json:"driver"`
	TripID  string      `json:"tripID"`
	RiderID string      `json:"riderID"`
}

type PaymentTripResponseData struct {
	TripID    string `json:"tripID"`
	UserID    string `json:"userID"`
	DriverID  string `json:"driverID"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	SessionID string `json:"sessionID"`
}

type PaymentEventSessionCreatedData struct {
	TripID    string  `json:"tripID"`
	SessionID string  `json:"sessionID"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type PaymentStatusUpdateData struct {
	TripID   string `json:"tripID"`
	UserID   string `json:"userID"`
	DriverID string `json:"driverID"`
}

type UserEventData struct {
	User *pbu.User `json:"user"`
}

type DriverNotifiedData struct {
	TripID  string `json:"tripID"`
	RiderID string `json:"riderID"`
}

type NoDriversFoundData struct {
	TripID  string `json:"tripID"`
	RiderID string `json:"riderID"`
}
