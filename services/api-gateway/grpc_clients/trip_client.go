package grpc_client

import (
	"os"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TripServiceClient wraps the gRPC client connection
type TripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*TripServiceClient, error) {
	tripServiceURL := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceURL == "" {
		tripServiceURL = "127.0.0.1:9093"
	}

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(tripServiceURL, dialOptions...)
	if err != nil {
		return nil, err
	}

	return &TripServiceClient{
		Client: pb.NewTripServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *TripServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
