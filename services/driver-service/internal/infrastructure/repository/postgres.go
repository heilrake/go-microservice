package repository

import (
	"context"
	"fmt"

	"ride-sharing/services/driver-service/internal/domain"
	db "ride-sharing/services/driver-service/internal/infrastructure/db"

	"github.com/google/uuid"
	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(gormDB *gorm.DB) domain.DriverRepository {
	return &postgresRepository{
		db: gormDB,
	}
}

func (r *postgresRepository) CreateDriver(ctx context.Context, userID, name, profilePicture string) (*db.DriverModel, error) {
	driver := &db.DriverModel{
		ID:             uuid.New().String(),
		UserID:         userID,
		Name:           name,
		ProfilePicture: profilePicture,
		IsAvailable:    false,
	}
	if err := r.db.WithContext(ctx).Create(driver).Error; err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}
	return driver, nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID string) (*db.DriverModel, error) {
	var d db.DriverModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *postgresRepository) CreateCar(ctx context.Context, driverID, carPlate, packageSlug string) (*db.CarModel, error) {
	car := &db.CarModel{
		ID:          uuid.New().String(),
		DriverID:    driverID,
		CarPlate:    carPlate,
		PackageSlug: packageSlug,
	}
	if err := r.db.WithContext(ctx).Create(car).Error; err != nil {
		return nil, fmt.Errorf("failed to create car: %w", err)
	}
	return car, nil
}

func (r *postgresRepository) ListCarsByDriverID(ctx context.Context, driverID string) ([]*db.CarModel, error) {
	var cars []*db.CarModel
	if err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).Find(&cars).Error; err != nil {
		return nil, fmt.Errorf("failed to list cars: %w", err)
	}
	return cars, nil
}

func (r *postgresRepository) GetCarByID(ctx context.Context, carID string) (*db.CarModel, error) {
	var car db.CarModel
	if err := r.db.WithContext(ctx).Where("id = ?", carID).First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *postgresRepository) RegisterDriver(ctx context.Context, userID, carID string, lat, lon float64) (*db.DriverModel, error) {
	driver, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	car, err := r.GetCarByID(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("car not found")
	}
	if car.DriverID != driver.ID {
		return nil, fmt.Errorf("car does not belong to driver")
	}
	gh := geohash.Encode(lat, lon)
	updates := map[string]interface{}{
		"is_available":   true,
		"current_car_id": carID,
		"latitude":       lat,
		"longitude":      lon,
		"geohash":        gh,
	}
	if err := r.db.WithContext(ctx).Model(&db.DriverModel{}).Where("id = ?", driver.ID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to register driver: %w", err)
	}
	driver.IsAvailable = true
	driver.CurrentCarID = &carID
	driver.Latitude = lat
	driver.Longitude = lon
	driver.Geohash = gh
	driver.CurrentCar = car
	return driver, nil
}

func (r *postgresRepository) UnregisterDriver(ctx context.Context, userID string) error {
	updates := map[string]interface{}{
		"is_available":   false,
		"current_car_id": nil,
		"latitude":       0,
		"longitude":      0,
		"geohash":        "",
	}
	if err := r.db.WithContext(ctx).Model(&db.DriverModel{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to unregister driver: %w", err)
	}
	return nil
}

func (r *postgresRepository) FindAvailableDrivers(ctx context.Context, packageType string) ([]*db.DriverModel, error) {
	var drivers []*db.DriverModel
	if err := r.db.WithContext(ctx).
		Preload("CurrentCar").
		Joins("INNER JOIN cars ON drivers.current_car_id = cars.id").
		Where("drivers.is_available = ? AND cars.package_slug = ?", true, packageType).
		Find(&drivers).Error; err != nil {
		return nil, fmt.Errorf("failed to find available drivers: %w", err)
	}
	return drivers, nil

}
