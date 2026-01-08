package infrastructure

import (
	"time"

	pb "ride-sharing/shared/proto/driver"
)

// Driver represents the drivers table in the database
type DriverModel struct {
	ID             string    `gorm:"column:id;type:uuid;primaryKey"`
	Name           string    `gorm:"column:name;type:varchar(255);not null"`
	ProfilePicture string    `gorm:"column:profile_picture;type:text"`
	CarPlate       string    `gorm:"column:car_plate;type:varchar(50)"`
	Geohash        string    `gorm:"column:geohash;type:varchar(20)"`
	PackageSlug    string    `gorm:"column:package_slug;type:varchar(50);not null"`
	Latitude       float64   `gorm:"column:latitude;type:decimal(10,8)"`
	Longitude      float64   `gorm:"column:longitude;type:decimal(11,8)"`
	IsAvailable    bool      `gorm:"column:is_available;type:boolean;default:true"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (DriverModel) TableName() string {
	return "drivers"
}

// ToDomain converts a database Driver to a proto Driver
func (d *DriverModel) ToProto() *pb.Driver {
	return &pb.Driver{
		Id:             d.ID,
		Name:           d.Name,
		ProfilePicture: d.ProfilePicture,
		CarPlate:       d.CarPlate,
		Geohash:        d.Geohash,
		PackageSlug:    d.PackageSlug,
		Location: &pb.Location{
			Latitude:  d.Latitude,
			Longitude: d.Longitude,
		},
	}
}

// FromProtoDriver converts a proto Driver to a database Driver
func FromProtoDriver(d *pb.Driver) *DriverModel {
	driver := &DriverModel{
		ID:             d.Id,
		Name:           d.Name,
		ProfilePicture: d.ProfilePicture,
		CarPlate:       d.CarPlate,
		Geohash:        d.Geohash,
		PackageSlug:    d.PackageSlug,
		IsAvailable:    true,
	}

	if d.Location != nil {
		driver.Latitude = d.Location.Latitude
		driver.Longitude = d.Location.Longitude
	}

	return driver
}
