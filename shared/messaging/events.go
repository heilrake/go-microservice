package messaging

import pb "ride-sharing/shared/proto/trip"

type TripEventData struct {
	Trip *pb.Trip `json:"trip"`
}
