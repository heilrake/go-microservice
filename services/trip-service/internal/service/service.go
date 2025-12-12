package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"github.com/google/uuid"
)

type TripService interface {
	CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error)
}

type tripService struct {
	repo repository.TripRepository
}

// NewTripServer wires a TripService with the given repository.
func NewTripServer(r repository.TripRepository) TripService {
	return &tripService{
		repo: r,
	}
}

func (s *tripService) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	trip := &domain.TripModel{
		ID:       uuid.NewString(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, trip)
}

func (s *tripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var routeResp tripTypes.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &routeResp, nil
}
