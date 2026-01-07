package infrastructure

import (
	"encoding/json"
	"time"

	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"
)

// Trip represents the trips table in the database
type Trip struct {
	ID                   string    `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID               string    `gorm:"column:user_id;type:varchar(255);not null"`
	Status               string    `gorm:"column:status;type:varchar(50);not null;default:'pending'"`
	RideFareID           *string   `gorm:"column:ride_fare_id;type:uuid"`
	DriverID             *string   `gorm:"column:driver_id;type:varchar(255)"`
	DriverName           *string   `gorm:"column:driver_name;type:varchar(255)"`
	DriverCarPlate       *string   `gorm:"column:driver_car_plate;type:varchar(50)"`
	DriverProfilePicture *string   `gorm:"column:driver_profile_picture;type:text"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Associations
	RideFare *RideFare `gorm:"foreignKey:RideFareID;references:ID"`
}

func (Trip) TableName() string {
	return "trips"
}

// RideFare represents the ride_fares table in the database
type RideFare struct {
	ID                string     `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID            string     `gorm:"column:user_id;type:varchar(255);not null"`
	PackageSlug       string     `gorm:"column:package_slug;type:varchar(50);not null"`
	TotalPriceInCents float64    `gorm:"column:total_price_in_cents;type:decimal(12,2);not null"`
	RouteData         []byte     `gorm:"column:route_data;type:jsonb"`
	ExpiresAt         *time.Time `gorm:"column:expires_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (RideFare) TableName() string {
	return "ride_fares"
}

// ToDomain converts a database Trip to a domain TripModel
func (t *Trip) ToDomain() *domain.TripModel {
	trip := &domain.TripModel{
		ID:     t.ID,
		UserID: t.UserID,
		Status: t.Status,
	}

	if t.RideFare != nil {
		trip.RideFare = t.RideFare.ToDomain()
	}

	if t.DriverID != nil {
		trip.Driver = &pb.TripDriver{
			Id: *t.DriverID,
		}
		if t.DriverName != nil {
			trip.Driver.Name = *t.DriverName
		}
		if t.DriverCarPlate != nil {
			trip.Driver.CarPlate = *t.DriverCarPlate
		}
		if t.DriverProfilePicture != nil {
			trip.Driver.ProfilePicture = *t.DriverProfilePicture
		}
	}

	return trip
}

// FromDomainTrip converts a domain TripModel to a database Trip
func FromDomainTrip(t *domain.TripModel) *Trip {
	trip := &Trip{
		ID:     t.ID,
		UserID: t.UserID,
		Status: t.Status,
	}

	if t.RideFare != nil {
		trip.RideFareID = &t.RideFare.ID
	}

	if t.Driver != nil {
		trip.DriverID = &t.Driver.Id
		trip.DriverName = &t.Driver.Name
		trip.DriverCarPlate = &t.Driver.CarPlate
		trip.DriverProfilePicture = &t.Driver.ProfilePicture
	}

	return trip
}

// ToDomain converts a database RideFare to a domain RideFareModel
func (r *RideFare) ToDomain() *domain.RideFareModel {
	fare := &domain.RideFareModel{
		ID:                r.ID,
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}

	if r.ExpiresAt != nil {
		fare.ExpiresAt = *r.ExpiresAt
	}

	if r.RouteData != nil {
		var route tripTypes.OsrmApiResponse
		if err := json.Unmarshal(r.RouteData, &route); err == nil {
			fare.Route = &route
		}
	}

	return fare
}

// FromDomainRideFare converts a domain RideFareModel to a database RideFare
func FromDomainRideFare(f *domain.RideFareModel) *RideFare {
	fare := &RideFare{
		ID:                f.ID,
		UserID:            f.UserID,
		PackageSlug:       f.PackageSlug,
		TotalPriceInCents: f.TotalPriceInCents,
	}

	if !f.ExpiresAt.IsZero() {
		fare.ExpiresAt = &f.ExpiresAt
	}

	if f.Route != nil {
		routeData, err := json.Marshal(f.Route)
		if err == nil {
			fare.RouteData = routeData
		}
	}

	return fare
}
