package domain

import pb "ride-sharing/shared/proto/trip"

type TripModel struct {
	ID       string
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}
