package service

import (
	"context"
	"fmt"

	"ride-sharing/services/user-service/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(r domain.UserRepository) domain.UserService {
	return &userService{
		repo: r,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, email, password, profilePicture string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:             uuid.New().String(),
		Username:       username,
		Email:          email,
		Password:       string(hashedPassword),
		ProfilePicture: profilePicture,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, username, email, profilePicture *string) (*domain.User, error) {
	user := &domain.User{ID: id}

	if username != nil {
		user.Username = *username
	}
	if email != nil {
		user.Email = *email
	}
	if profilePicture != nil {
		user.ProfilePicture = *profilePicture
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Fetch updated user
	return s.repo.GetByID(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
