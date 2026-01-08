package repository

import (
	"context"
	"fmt"

	math "math/rand/v2"
	"ride-sharing/shared/util"

	"ride-sharing/services/driver-service/internal/domain"
	db "ride-sharing/services/driver-service/internal/infrastructure/db"
	pb "ride-sharing/shared/proto/driver"

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

func (r *postgresRepository) CreateDriver(ctx context.Context, packageSlug string) (*db.DriverModel, error) {

	randomIndex := math.IntN(len(util.PredefinedRoutes))
	randomRoute := util.PredefinedRoutes[randomIndex]

	randomPlate := util.GenerateRandomPlate()
	randomAvatar := util.GetRandomAvatar(randomIndex)

	// we can ignore this property for now, but it must be sent to the frontend.
	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	dbDriver := db.FromProtoDriver(&pb.Driver{
		PackageSlug:    packageSlug,
		Id:             uuid.New().String(),
		Name:           "Lando Norris",
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlate,
		Geohash:        geohash,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
	})

	if err := r.db.WithContext(ctx).Create(dbDriver).Error; err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	return dbDriver, nil
}

func (r *postgresRepository) UnregisterDriver(ctx context.Context, driverID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", driverID).Delete(&db.DriverModel{}).Error; err != nil {
		return fmt.Errorf("failed to unregister driver: %w", err)
	}

	return nil
}

func (r *postgresRepository) FindAvailableDrivers(ctx context.Context, packageType string) ([]*db.DriverModel, error) {
	var drivers []*db.DriverModel

	if err := r.db.WithContext(ctx).Where("package_slug = ? AND is_available = ?", packageType, true).Find(&drivers).Error; err != nil {
		return nil, fmt.Errorf("failed to find available drivers: %w", err)
	}

	return drivers, nil
}
