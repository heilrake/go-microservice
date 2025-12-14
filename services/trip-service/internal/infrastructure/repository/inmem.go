package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)

// inmemRepository provides an in-memory implementation of TripRepository.
// Safe for single-threaded demo/testing purposes only.
type inmemRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error)
	SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error
	GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error)
}

// NewInmemRepository creates a new in-memory repository instance.
func NewInmemRepository() TripRepository {
	return &inmemRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inmemRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID] = trip

	return trip, nil
}

func (r *inmemRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.rideFares[fare.ID] = fare

	return nil
}

func (r *inmemRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) {
	fare, exists := r.rideFares[id]
	if !exists {
		return nil, nil
	}

	return fare, nil
}
