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

func (h *driverGrpcHandler) CreateDriver(ctx context.Context, req *pb.CreateDriverRequest) (*pb.CreateDriverResponse, error) {
	driver, err := h.service.CreateDriver(ctx, req.GetUserId(), req.GetName(), req.GetProfilePicture())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create driver: %v", err)
	}
	return &pb.CreateDriverResponse{Driver: driver.ToProto()}, nil
}

func (h *driverGrpcHandler) CreateCar(ctx context.Context, req *pb.CreateCarRequest) (*pb.CreateCarResponse, error) {
	car, err := h.service.CreateCar(ctx, req.GetUserId(), req.GetCarPlate(), req.GetPackageSlug())
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &pb.CreateCarResponse{Car: car.ToProtoCar()}, nil
}

func (h *driverGrpcHandler) ListCars(ctx context.Context, req *pb.ListCarsRequest) (*pb.ListCarsResponse, error) {
	cars, err := h.service.ListCars(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	out := make([]*pb.Car, len(cars))
	for i, c := range cars {
		out[i] = c.ToProtoCar()
	}
	return &pb.ListCarsResponse{Cars: out}, nil
}

func (h *driverGrpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driver, err := h.service.RegisterDriver(ctx, req.GetDriverID(), req.GetCarId(), req.GetLatitude(), req.GetLongitude())
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "%s", err.Error())
	}
	return &pb.RegisterDriverResponse{Driver: driver.ToProto()}, nil
}

func (h *driverGrpcHandler) UnRegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	if err := h.service.UnregisterDriver(ctx, req.GetDriverID()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unregister driver")
	}
	return &pb.RegisterDriverResponse{Driver: &pb.Driver{Id: req.GetDriverID()}}, nil
}
