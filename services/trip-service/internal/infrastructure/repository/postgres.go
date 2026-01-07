package repository

import (
	"context"
	"fmt"

	"ride-sharing/services/trip-service/internal/domain"
	db "ride-sharing/services/trip-service/internal/infrastructure/db"

	pbd "ride-sharing/shared/proto/driver"

	"gorm.io/gorm"
)

// postgresRepository provides a PostgreSQL implementation of TripRepository using GORM.
type postgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgreSQL repository instance.
func NewPostgresRepository(gormDB *gorm.DB) TripRepository {
	return &postgresRepository{
		db: gormDB,
	}
}

func (r *postgresRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	dbTrip := db.FromDomainTrip(trip)

	if err := r.db.WithContext(ctx).Create(dbTrip).Error; err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	result := dbTrip.ToDomain()
	result.RideFare = trip.RideFare
	return result, nil
}

func (r *postgresRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	dbFare := db.FromDomainRideFare(fare)

	if err := r.db.WithContext(ctx).Create(dbFare).Error; err != nil {
		return fmt.Errorf("failed to save ride fare: %w", err)
	}

	return nil
}

func (r *postgresRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	var dbFare db.RideFare

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbFare).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ride fare: %w", err)
	}

	return dbFare.ToDomain(), nil
}

func (r *postgresRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	var dbTrip db.Trip

	if err := r.db.WithContext(ctx).Preload("RideFare").Where("id = ?", id).First(&dbTrip).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	return dbTrip.ToDomain(), nil
}

func (r *postgresRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if driver != nil {
		updates["driver_id"] = driver.Id
		updates["driver_name"] = driver.Name
		updates["driver_car_plate"] = driver.CarPlate
		updates["driver_profile_picture"] = driver.ProfilePicture
	}

	result := r.db.WithContext(ctx).Model(&db.Trip{}).Where("id = ?", tripID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update trip: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("trip not found with ID: %s", tripID)
	}

	return nil
}
