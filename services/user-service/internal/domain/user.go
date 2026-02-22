package domain

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type UserService interface {
	CreateUser(ctx context.Context, username, email, password, profilePicture, role string) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	LoginUser(ctx context.Context, email, password, role string) (*User, error)
	GetOrCreateUserByOAuth(ctx context.Context, email, username, profilePicture, role string) (*User, error)
	UpdateUser(ctx context.Context, id string, username, email, profilePicture *string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

// User is the domain model for a user
type User struct {
	ID             string
	Username       string
	Email          string
	Password       string
	ProfilePicture string
	Role           string // "rider" | "driver"
}
