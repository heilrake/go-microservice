package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"ride-sharing/services/payment-service/internal/domain"
)

type paymentService struct {
	repo             domain.Repository
	paymentProcessor domain.PaymentProcessor
}

func NewPaymentService(repo domain.Repository, paymentProcessor domain.PaymentProcessor) domain.Service {
	return &paymentService{
		repo:             repo,
		paymentProcessor: paymentProcessor,
	}
}

func (s *paymentService) CreatePaymentIntent(ctx context.Context, tripID, userID string, amount int64, currency string) (*domain.PaymentIntentModel, error) {
	metadata := map[string]string{
		"trip_id": tripID,
		"user_id": userID,
	}

	stripeID, clientSecret, err := s.paymentProcessor.CreatePaymentIntent(ctx, amount, currency, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	intent := &domain.PaymentIntentModel{
		ID:                    uuid.New().String(),
		TripID:                tripID,
		UserID:                userID,
		StripePaymentIntentID: stripeID,
		ClientSecret:          clientSecret,
		Amount:                amount,
		Currency:              currency,
		Status:                "authorized",
	}

	if err := s.repo.Save(ctx, intent); err != nil {
		return nil, fmt.Errorf("failed to save payment intent: %w", err)
	}

	return intent, nil
}

func (s *paymentService) CapturePayment(ctx context.Context, tripID string) error {
	intent, err := s.repo.GetByTripID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to get payment intent: %w", err)
	}

	if err := s.paymentProcessor.CapturePayment(ctx, intent.StripePaymentIntentID, nil); err != nil {
		return fmt.Errorf("failed to capture payment: %w", err)
	}

	return s.repo.UpdateStatus(ctx, tripID, "captured")
}

func (s *paymentService) CancelPayment(ctx context.Context, tripID string) error {
	intent, err := s.repo.GetByTripID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed to get payment intent: %w", err)
	}

	if err := s.paymentProcessor.CancelPayment(ctx, intent.StripePaymentIntentID); err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	return s.repo.UpdateStatus(ctx, tripID, "cancelled")
}
