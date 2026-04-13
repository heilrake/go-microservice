package stripe

import (
	"context"
	"fmt"

	stripeApi "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"

	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/services/payment-service/pkg/types"
)

type stripeClient struct {
	config *types.PaymentConfig
}

func NewStripeClient(config *types.PaymentConfig) domain.PaymentProcessor {
	stripeApi.Key = config.StripeSecretKey
	return &stripeClient{config: config}
}

func (s *stripeClient) CreatePaymentIntent(ctx context.Context, amount int64, currency string, metadata map[string]string) (stripeID, clientSecret string, err error) {
	params := &stripeApi.PaymentIntentParams{
		Amount:        stripeApi.Int64(amount),
		Currency:      stripeApi.String(currency),
		CaptureMethod: stripeApi.String(string(stripeApi.PaymentIntentCaptureMethodManual)),
		Metadata:      metadata,
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", "", fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi.ID, pi.ClientSecret, nil
}

func (s *stripeClient) CapturePayment(ctx context.Context, paymentIntentID string, amountToCapture *int64) error {
	params := &stripeApi.PaymentIntentCaptureParams{}

	if amountToCapture != nil {
		params.AmountToCapture = stripeApi.Int64(*amountToCapture)
	}

	_, err := paymentintent.Capture(paymentIntentID, params)
	if err != nil {
		return fmt.Errorf("failed to capture payment: %w", err)
	}

	return nil
}

func (s *stripeClient) CancelPayment(ctx context.Context, paymentIntentID string) error {
	_, err := paymentintent.Cancel(paymentIntentID, nil)
	if err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	return nil
}
