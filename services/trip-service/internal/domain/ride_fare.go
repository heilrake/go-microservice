package domain

import (
	"ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"
	"time"
)

type RideFareModel struct {
	ID                string
	UserID            string
	PackageSlug       string // e.g., "standard", "premium"
	TotalPriceInCents float64
	Route             *types.OsrmApiResponse
	ExpiresAt         time.Time
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID,
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare
	for _, f := range fares {
		protoFares = append(protoFares, f.ToProto())
	}
	return protoFares
}
