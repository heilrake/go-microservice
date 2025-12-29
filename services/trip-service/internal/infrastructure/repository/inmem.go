package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"

	pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
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
	GetTripByID(ctx context.Context, id string) (*domain.TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error
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

func (r *inmemRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) {
	trip, exists := r.trips[id]
	if !exists {
		return nil, nil
	}

	return trip, nil
}

func (r *inmemRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error {
	trip, ok := r.trips[tripID]
	if !ok {
		return fmt.Errorf("trip not found with ID: %s", tripID)
	}

	trip.Status = status

	if driver != nil {
		trip.Driver = &pb.TripDriver{
			Id:             driver.Id,
			Name:           driver.Name,
			CarPlate:       driver.CarPlate,
			ProfilePicture: driver.ProfilePicture,
		}
	}
	return nil
}
