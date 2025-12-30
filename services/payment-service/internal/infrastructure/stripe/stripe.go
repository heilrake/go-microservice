package stripe

import (
	"context"
	"fmt"
	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/services/payment-service/pkg/types"

	stripeApi "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type stripeClient struct {
	config *types.PaymentConfig
}

func NewStripeClient(config *types.PaymentConfig) domain.PaymentProcessor {

	stripeApi.Key = config.StripeSecretKey

	return &stripeClient{
		config: config,
	}
}

func (s *stripeClient) CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error) {

	params := &stripeApi.CheckoutSessionParams{
		SuccessURL: stripeApi.String(s.config.SuccessURL),
		CancelURL:  stripeApi.String(s.config.CancelURL),
		Metadata:   metadata,
		LineItems: []*stripeApi.CheckoutSessionLineItemParams{
			{
				PriceData: &stripeApi.CheckoutSessionLineItemPriceDataParams{
					Currency: stripeApi.String(currency),
					ProductData: &stripeApi.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripeApi.String("Ride Payment"),
					},
					UnitAmount: stripeApi.Int64(amount),
				},
				Quantity: stripeApi.Int64(1),
			},
		},
		Mode: stripeApi.String(string(stripeApi.CheckoutSessionModePayment)),
	}

	result, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create a payment session on stripe: %w", err)
	}

	return result.ID, nil
}
