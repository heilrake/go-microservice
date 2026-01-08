package grpc

import (
	"context"
	"ride-sharing/services/driver-service/internal/domain"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type driverGrpcHandler struct {
	pb.UnimplementedDriverServiceServer

	service domain.DriverService
}

func NewGrpcHandler(s *grpc.Server, service domain.DriverService) {
	handler := &driverGrpcHandler{
		service: service,
	}

	pb.RegisterDriverServiceServer(s, handler)
}

func (h *driverGrpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driver, err := h.service.RegisterDriver(ctx, req.GetPackageSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver")
	}

	return &pb.RegisterDriverResponse{
		Driver: driver.ToProto(),
	}, nil
}

func (h *driverGrpcHandler) UnRegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {

	if err := h.service.UnregisterDriver(ctx, req.GetDriverID()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unregister driver")
	}

	return &pb.RegisterDriverResponse{
		Driver: &pb.Driver{
			Id: req.GetDriverID(),
		},
	}, nil
}
