package domain

import (
	"context"

	pb "ride-sharing/shared/proto/driver"
)

type Driver struct {
	ID             string
	UserID         string
	Name           string
	ProfilePicture string
	Geohash        string
	Latitude       float64
	Longitude      float64
	CurrentCarID   *string
	IsAvailable    bool
	CurrentCar     *Car
}

type Car struct {
	ID          string
	DriverID    string
	CarPlate    string
	PackageSlug string
}

func (d *Driver) ToProto() *pb.Driver {
	driver := &pb.Driver{
		Id:             d.ID,
		UserId:         d.UserID,
		Name:           d.Name,
		ProfilePicture: d.ProfilePicture,
		Geohash:        d.Geohash,
		Location: &pb.Location{
			Latitude:  d.Latitude,
			Longitude: d.Longitude,
		},
	}
	if d.CurrentCar != nil {
		driver.CarPlate = d.CurrentCar.CarPlate
		driver.PackageSlug = d.CurrentCar.PackageSlug
	}
	return driver
}

func (c *Car) ToProtoCar() *pb.Car {
	return &pb.Car{
		Id:          c.ID,
		DriverId:    c.DriverID,
		CarPlate:    c.CarPlate,
		PackageSlug: c.PackageSlug,
	}
}

type DriverRepository interface {
	CreateDriver(ctx context.Context, userID, name, profilePicture string) (*Driver, error)
	GetByUserID(ctx context.Context, userID string) (*Driver, error)
	CreateCar(ctx context.Context, driverID, carPlate, packageSlug string) (*Car, error)
	ListCarsByDriverID(ctx context.Context, driverID string) ([]*Car, error)
	GetCarByID(ctx context.Context, carID string) (*Car, error)
	RegisterDriver(ctx context.Context, userID, carID string, lat, lon float64) (*Driver, error)
	UnregisterDriver(ctx context.Context, userID string) error
	FindAvailableDrivers(ctx context.Context, packageType string) ([]*Driver, error)
}

type DriverService interface {
	CreateDriver(ctx context.Context, userID, name, profilePicture string) (*Driver, error)
	GetDriver(ctx context.Context, userID string) (*Driver, error)
	CreateCar(ctx context.Context, userID, carPlate, packageSlug string) (*Car, error)
	ListCars(ctx context.Context, userID string) ([]*Car, error)
	RegisterDriver(ctx context.Context, userID, carID string, lat, lon float64) (*Driver, error)
	UnregisterDriver(ctx context.Context, userID string) error
	FindAvailableDrivers(ctx context.Context, packageType string) ([]*Driver, error)
}
