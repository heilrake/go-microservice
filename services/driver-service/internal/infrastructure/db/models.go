package infrastructure

import (
	"time"

	pb "ride-sharing/shared/proto/driver"
)

// DriverModel represents the drivers table
type DriverModel struct {
	ID             string    `gorm:"column:id;type:uuid;primaryKey"`
	UserID         string    `gorm:"column:user_id;type:uuid;uniqueIndex"`
	Name           string    `gorm:"column:name;type:varchar(255);not null"`
	ProfilePicture string    `gorm:"column:profile_picture;type:text"`
	Geohash        string    `gorm:"column:geohash;type:varchar(20)"`
	Latitude       float64   `gorm:"column:latitude;type:decimal(10,8)"`
	Longitude      float64   `gorm:"column:longitude;type:decimal(11,8)"`
	CurrentCarID   *string   `gorm:"column:current_car_id;type:uuid"`
	IsAvailable    bool      `gorm:"column:is_available;type:boolean;default:true"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`

	CurrentCar *CarModel `gorm:"foreignKey:CurrentCarID"`
}

// CarModel represents the cars table
type CarModel struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey"`
	DriverID    string    `gorm:"column:driver_id;type:uuid;not null"`
	CarPlate    string    `gorm:"column:car_plate;type:varchar(50);not null"`
	PackageSlug string    `gorm:"column:package_slug;type:varchar(50);not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (CarModel) TableName() string {
	return "cars"
}

func (DriverModel) TableName() string {
	return "drivers"
}

// ToProto converts DriverModel to proto Driver (carPlate/packageSlug from current car if available)
func (d *DriverModel) ToProto() *pb.Driver {
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

// ToProtoCar converts CarModel to proto Car
func (c *CarModel) ToProtoCar() *pb.Car {
	return &pb.Car{
		Id:          c.ID,
		DriverId:    c.DriverID,
		CarPlate:    c.CarPlate,
		PackageSlug: c.PackageSlug,
	}
}
