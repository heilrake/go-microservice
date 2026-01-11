package main

import (
	grpc_client "ride-sharing/services/api-gateway/grpc_clients"
)

type Clients struct {
	UserClient   *grpc_client.UserServiceClient
	TripClient   *grpc_client.TripServiceClient
	DriverClient *grpc_client.DriverServiceClient
}

// NewApp initializes all gRPC clients once at startup
func InitializeClients() (*Clients, error) {
	userClient, err := grpc_client.NewUserServiceClient()
	if err != nil {
		return nil, err
	}

	tripClient, err := grpc_client.NewTripServiceClient()
	if err != nil {
		userClient.Close()
		return nil, err
	}

	driverClient, err := grpc_client.NewDriverServiceClient()
	if err != nil {
		userClient.Close()
		tripClient.Close()
		return nil, err
	}

	return &Clients{
		UserClient:   userClient,
		TripClient:   tripClient,
		DriverClient: driverClient,
	}, nil
}

func (a *Clients) CloseAllClients() {
	if a.UserClient != nil {
		a.UserClient.Close()
	}
	if a.TripClient != nil {
		a.TripClient.Close()
	}
	if a.DriverClient != nil {
		a.DriverClient.Close()
	}
}
