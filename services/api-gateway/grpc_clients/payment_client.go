package grpc_client

import (
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "ride-sharing/shared/proto/payment"
	"ride-sharing/shared/tracing"
)

type PaymentServiceClient struct {
	Client pb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentServiceClient() (*PaymentServiceClient, error) {
	addr := os.Getenv("PAYMENT_SERVICE_URL")
	if addr == "" {
		addr = "127.0.0.1:9004"
	}

	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.NewClient(addr, dialOptions...)
	if err != nil {
		return nil, err
	}

	return &PaymentServiceClient{
		Client: pb.NewPaymentServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *PaymentServiceClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
