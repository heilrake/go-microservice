package domain

import (
	"context"

	tripTypes "ride-sharing/services/trip-service/pkg/types"
	pbd "ride-sharing/shared/proto/driver"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type TripModel struct {
	ID       string
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

func (t *TripModel) ToProto() *pb.Trip {
	var route *pb.Route
	var selectedFare *pb.RideFare

	if t.RideFare != nil {
		selectedFare = t.RideFare.ToProto()
		if t.RideFare.Route != nil {
			route = t.RideFare.Route.ToProto()
		}
	}

	return &pb.Trip{
		Id:           t.ID,
		UserID:       t.UserID,
		SelectedFare: selectedFare,
		Status:       t.Status,
		Driver:       t.Driver,
		Route:        route,
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	GetTripByID(ctx context.Context, id string) (*TripModel, error)
	GetRideFareByID(ctx context.Context, id string) (*RideFareModel, error)
	SaveRideFare(ctx context.Context, fare *RideFareModel) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userID string, route *tripTypes.OsrmApiResponse) ([]*RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error
}
