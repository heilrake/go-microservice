package service

import (
	"context"
	"fmt"

	"ride-sharing/services/driver-service/internal/domain"
	infrastructure "ride-sharing/services/driver-service/internal/infrastructure/db"
)

type driverService struct {
	repo domain.DriverRepository
}

func NewDriverService(r domain.DriverRepository) domain.DriverService {
	return &driverService{
		repo: r,
	}
}

func (s *driverService) CreateDriver(ctx context.Context, userID, name, profilePicture string) (*infrastructure.DriverModel, error) {
	return s.repo.CreateDriver(ctx, userID, name, profilePicture)
}

func (s *driverService) GetDriver(ctx context.Context, userID string) (*infrastructure.DriverModel, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *driverService) CreateCar(ctx context.Context, userID, carPlate, packageSlug string) (*infrastructure.CarModel, error) {
	driver, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("driver not found: create driver profile first")
	}
	return s.repo.CreateCar(ctx, driver.ID, carPlate, packageSlug)
}

func (s *driverService) ListCars(ctx context.Context, userID string) ([]*infrastructure.CarModel, error) {
	driver, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("driver not found")
	}
	return s.repo.ListCarsByDriverID(ctx, driver.ID)
}

func (s *driverService) RegisterDriver(ctx context.Context, userID, carID string, lat, lon float64) (*infrastructure.DriverModel, error) {
	return s.repo.RegisterDriver(ctx, userID, carID, lat, lon)
}

func (s *driverService) UnregisterDriver(ctx context.Context, userID string) error {
	return s.repo.UnregisterDriver(ctx, userID)
}

func (s *driverService) FindAvailableDrivers(ctx context.Context, packageType string) ([]*infrastructure.DriverModel, error) {
	return s.repo.FindAvailableDrivers(ctx, packageType)
}
