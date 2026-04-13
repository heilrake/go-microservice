package repository

import (
	"context"

	"gorm.io/gorm"

	"ride-sharing/services/payment-service/internal/domain"
	dbmodels "ride-sharing/services/payment-service/internal/infrastructure/db"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(gormDB *gorm.DB) domain.Repository {
	return &postgresRepository{db: gormDB}
}

func (r *postgresRepository) Save(ctx context.Context, intent *domain.PaymentIntentModel) error {
	record := &dbmodels.PaymentIntent{
		ID:                    intent.ID,
		TripID:                intent.TripID,
		UserID:                intent.UserID,
		StripePaymentIntentID: intent.StripePaymentIntentID,
		Amount:                intent.Amount,
		Currency:              intent.Currency,
		Status:                intent.Status,
	}
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *postgresRepository) GetByTripID(ctx context.Context, tripID string) (*domain.PaymentIntentModel, error) {
	var record dbmodels.PaymentIntent
	if err := r.db.WithContext(ctx).Where("trip_id = ?", tripID).First(&record).Error; err != nil {
		return nil, err
	}
	return &domain.PaymentIntentModel{
		ID:                    record.ID,
		TripID:                record.TripID,
		UserID:                record.UserID,
		StripePaymentIntentID: record.StripePaymentIntentID,
		Amount:                record.Amount,
		Currency:              record.Currency,
		Status:                record.Status,
	}, nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, tripID string, status string) error {
	return r.db.WithContext(ctx).
		Model(&dbmodels.PaymentIntent{}).
		Where("trip_id = ?", tripID).
		Update("status", status).Error
}
