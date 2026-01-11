package grpc_client

import (
	"os"
	pb "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DriverServiceClient wraps the gRPC client connection
type DriverServiceClient struct {
	Client pb.DriverServiceClient
	conn   *grpc.ClientConn
}

func NewDriverServiceClient() (*DriverServiceClient, error) {
	driverServiceURL := os.Getenv("DRIVER_SERVICE_URL")
	if driverServiceURL == "" {
		driverServiceURL = "127.0.0.1:9094"
	}

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(driverServiceURL, dialOptions...)
	if err != nil {
		return nil, err
	}

	return &DriverServiceClient{
		Client: pb.NewDriverServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *DriverServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
