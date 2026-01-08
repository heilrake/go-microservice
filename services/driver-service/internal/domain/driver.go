package domain

import (
	"context"
	infrastructure "ride-sharing/services/driver-service/internal/infrastructure/db"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, packageSlug string) (*infrastructure.DriverModel, error)
	UnregisterDriver(ctx context.Context, driverID string) error
	FindAvailableDrivers(ctx context.Context, packageType string) ([]*infrastructure.DriverModel, error)
}

type DriverService interface {
	RegisterDriver(ctx context.Context, packageSlug string) (*infrastructure.DriverModel, error)
	UnregisterDriver(ctx context.Context, driverID string) error
	FindAvailableDrivers(ctx context.Context, packageType string) ([]*infrastructure.DriverModel, error)
}
