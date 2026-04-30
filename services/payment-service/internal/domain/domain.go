package domain

import "context"

const (
	StatusAuthorized = "authorized"
	StatusCaptured   = "captured"
	StatusCancelled  = "cancelled"
)

type Service interface {
	CreatePaymentIntent(ctx context.Context, tripID, userID string, amount int64, currency string) (*PaymentIntentModel, error)
	CapturePayment(ctx context.Context, tripID string) error
	CancelPayment(ctx context.Context, tripID string) error
}

type Repository interface {
	Save(ctx context.Context, intent *PaymentIntentModel) error
	GetByTripID(ctx context.Context, tripID string) (*PaymentIntentModel, error)
	UpdateStatus(ctx context.Context, tripID string, status string) error
}

type PaymentProcessor interface {
	CreatePaymentIntent(ctx context.Context, amount int64, currency string, metadata map[string]string) (stripeID, clientSecret string, err error)
	CapturePayment(ctx context.Context, paymentIntentID string, amountToCapture *int64) error
	CancelPayment(ctx context.Context, paymentIntentID string) error
}

type PaymentIntentModel struct {
	ID                    string
	TripID                string
	UserID                string
	StripePaymentIntentID string
	ClientSecret          string
	Amount                int64
	Currency              string
	Status                string
}
