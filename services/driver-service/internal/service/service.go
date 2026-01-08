package service

import (
	"context"
	"fmt"
	"ride-sharing/services/driver-service/internal/domain"
	infrastructure "ride-sharing/services/driver-service/internal/infrastructure/db"
	pb "ride-sharing/shared/proto/driver"
)

type driverInMap struct {
	Driver *pb.Driver
	// Index int
	// TODO: route
}
type driverService struct {
	repo domain.DriverRepository
}

func NewDriverService(r domain.DriverRepository) domain.DriverService {
	return &driverService{
		repo: r,
	}
}

func (s *driverService) RegisterDriver(ctx context.Context, packageSlug string) (*infrastructure.DriverModel, error) {
	driver, err := s.repo.CreateDriver(ctx, packageSlug)
	if err != nil {
		return nil, fmt.Errorf("failed to register driver: %w", err)
	}

	return driver, nil
}

func (s *driverService) UnregisterDriver(ctx context.Context, driverID string) error {
	return s.repo.UnregisterDriver(ctx, driverID)
}

func (s *driverService) FindAvailableDrivers(ctx context.Context, packageType string) ([]*infrastructure.DriverModel, error) {
	drivers, err := s.repo.FindAvailableDrivers(ctx, packageType)
	if err != nil {
		return nil, fmt.Errorf("failed to find available drivers: %w", err)
	}

	return drivers, nil
}
