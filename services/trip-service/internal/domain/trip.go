package domain

import pb "ride-sharing/shared/proto/trip"

type TripModel struct {
	ID       string
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

func (t *TripModel) ToProto() *pb.Trip {
	var route *pb.Route
	if t.RideFare != nil && t.RideFare.Route != nil {
		route = t.RideFare.Route.ToProto()
	}

	return &pb.Trip{
		Id:           t.ID,
		UserID:       t.UserID,
		SelectedFare: t.RideFare.ToProto(),
		Status:       t.Status,
		Driver:       t.Driver,
		Route:        route,
	}
}
