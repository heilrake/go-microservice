package main

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type previewTripRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (p *previewTripRequest) toProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: p.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  p.Pickup.Latitude,
			Longitude: p.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  p.Destination.Latitude,
			Longitude: p.Destination.Longitude,
		},
	}
}

type startTripRequest struct {
	RideFareID string `json:"rideFareID"`
	UserID     string `json:"userID"`
}

type cancelTripRequest struct {
	UserID string `json:"userID"`
}

func (c *cancelTripRequest) toProto() *pb.CancelTripRequest {
	return &pb.CancelTripRequest{
		UserID: c.UserID,
	}
}

func (c *startTripRequest) toProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		RideFareID: c.RideFareID,
		UserID:     c.UserID,
	}
}

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // "rider" | "driver" — required
}

type createDriverRequest struct {
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
}

type createCarRequest struct {
	CarPlate    string `json:"car_plate"`
	PackageSlug string `json:"package_slug"`
}

// oauthLoginRequest is the request body for POST /auth/oauth.
type oauthLoginRequest struct {
	Code        string `json:"code"`
	Provider    string `json:"provider"` // "google"
	Role        string `json:"role"`     // "rider" | "driver"
	RedirectURI string `json:"redirect_uri"`
}

// devLoginRequest is the request body for POST /dev/login.
type devLoginRequest struct {
	Role string `json:"role"` // "rider" | "driver"
	Seed int    `json:"seed"`
}
