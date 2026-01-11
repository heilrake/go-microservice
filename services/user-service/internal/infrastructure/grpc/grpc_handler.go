package grpc

import (
	"context"

	"ride-sharing/services/user-service/internal/domain"
	"ride-sharing/services/user-service/internal/infrastructure/events"
	pb "ride-sharing/shared/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userGrpcHandler struct {
	pb.UnimplementedUserServiceServer

	service   domain.UserService
	publisher *events.UserEventPublisher
}

func NewGrpcHandler(s *grpc.Server, svc domain.UserService, publisher *events.UserEventPublisher) {
	handler := &userGrpcHandler{
		service:   svc,
		publisher: publisher,
	}

	pb.RegisterUserServiceServer(s, handler)
}

func (h *userGrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.service.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
		},
	}, nil
}

func (h *userGrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := h.service.CreateUser(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword(), req.GetProfilePicture())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Publish user created event
	if err := h.publisher.PublishUserCreated(ctx, user); err != nil {
		// Log error but don't fail the request
		// In production, you might want to use outbox pattern for guaranteed delivery
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
		},
	}, nil
}

func (h *userGrpcHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	var username, email, profilePicture *string

	if req.Username != nil {
		username = req.Username
	}
	if req.Email != nil {
		email = req.Email
	}
	if req.ProfilePicture != nil {
		profilePicture = req.ProfilePicture
	}

	user, err := h.service.UpdateUser(ctx, req.GetUserId(), username, email, profilePicture)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// Publish user updated event - this will notify other services
	if err := h.publisher.PublishUserUpdated(ctx, user); err != nil {
		// Log error but don't fail the request
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			ProfilePicture: user.ProfilePicture,
		},
	}, nil
}
