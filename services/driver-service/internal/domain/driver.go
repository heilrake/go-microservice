package domain

import (
	"context"
	infrastructure "ride-sharing/services/driver-service/internal/infrastructure/db"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, userID, name, profilePicture string) (*infrastructure.DriverModel, error)
	GetByUserID(ctx context.Context, userID string) (*infrastructure.DriverModel, error)
	CreateCar(ctx context.Context, driverID, carPlate, packageSlug string) (*infrastructure.CarModel, error)
	ListCarsByDriverID(ctx context.Context, driverID string) ([]*infrastructure.CarModel, error)
	GetCarByID(ctx context.Context, carID string) (*infrastructure.CarModel, error)
	RegisterDriver(ctx context.Context, userID, carID string, lat, lon float64) (*infrastructure.DriverModel, error)
	UnregisterDriver(ctx context.Context, userID string) error
	FindAvailableDrivers(ctx context.Context, packageType string) ([]*infrastructure.DriverModel, error)
}

type DriverService interface {
	CreateDriver(ctx context.Context, userID, name, profilePicture string) (*infrastructure.DriverModel, error)
	CreateCar(ctx context.Context, userID, carPlate, packageSlug string) (*infrastructure.CarModel, error)
	ListCars(ctx context.Context, userID string) ([]*infrastructure.CarModel, error)
	RegisterDriver(ctx context.Context, userID, carID string, lat, lon float64) (*infrastructure.DriverModel, error)
	UnregisterDriver(ctx context.Context, userID string) error
	FindAvailableDrivers(ctx context.Context, packageType string) ([]*infrastructure.DriverModel, error)
}
