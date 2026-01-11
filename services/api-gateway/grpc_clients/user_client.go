package grpc_client

import (
	"os"
	pb "ride-sharing/shared/proto/user"
	"ride-sharing/shared/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userServiceClient struct {
	Client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserServiceClient() (*userServiceClient, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "localhost:9095"
	}

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(userServiceURL, dialOptions...)
	if err != nil {
		return nil, err
	}

	return &userServiceClient{
		Client: pb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *userServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
