package repository

import (
	"context"
	"fmt"

	"ride-sharing/services/user-service/internal/domain"
	db "ride-sharing/services/user-service/internal/infrastructure/db"

	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(gormDB *gorm.DB) domain.UserRepository {
	return &postgresRepository{
		db: gormDB,
	}
}

func (r *postgresRepository) Create(ctx context.Context, user *domain.User) error {
	dbUser := &db.UserModel{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Password:       user.Password,
		ProfilePicture: user.ProfilePicture,
	}

	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var dbUser db.UserModel

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbUser).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &domain.User{
		ID:             dbUser.ID,
		Username:       dbUser.Username,
		Email:          dbUser.Email,
		Password:       dbUser.Password,
		ProfilePicture: dbUser.ProfilePicture,
	}, nil
}

func (r *postgresRepository) Update(ctx context.Context, user *domain.User) error {
	updates := map[string]interface{}{}

	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.ProfilePicture != "" {
		updates["profile_picture"] = user.ProfilePicture
	}

	if len(updates) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Model(&db.UserModel{}).Where("id = ?", user.ID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&db.UserModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
