package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ride-sharing/services/payment-service/internal/domain"
	pb "ride-sharing/shared/proto/payment"
)

type grpcHandler struct {
	pb.UnimplementedPaymentServiceServer
	service domain.Service
}

func NewGRPCHandler(server *grpc.Server, service domain.Service) {
	handler := &grpcHandler{service: service}
	pb.RegisterPaymentServiceServer(server, handler)
}

func (h *grpcHandler) CreatePaymentIntent(ctx context.Context, req *pb.CreatePaymentIntentRequest) (*pb.CreatePaymentIntentResponse, error) {
	intent, err := h.service.CreatePaymentIntent(
		ctx,
		req.GetTripID(),
		req.GetUserID(),
		req.GetAmountInCents(),
		req.GetCurrency(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment intent: %v", err)
	}

	return &pb.CreatePaymentIntentResponse{
		PaymentIntentID: intent.StripePaymentIntentID,
		ClientSecret:    intent.ClientSecret,
	}, nil
}
