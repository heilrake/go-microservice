package internal

import (
	"os"
	pb "ride-sharing/shared/proto/user"
	"ride-sharing/shared/tracing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	Client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient() (*UserClient, error) {
	addr := os.Getenv("USER_SERVICE_URL")
	if addr == "" {
		addr = "127.0.0.1:9095"
	}
	opts := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}
	return &UserClient{
		Client: pb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
