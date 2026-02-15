package infrastructure

import (
	"time"

	"github.com/google/uuid"
)

type AuthUser struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex"`
	PasswordHash string
	Role         string // "passenger", "driver", "admin"
	IsActive     bool
	CreatedAt    time.Time
}

func (AuthUser) TableName() string {
	return "auth_users"
}
