package infrastructure

import (
	"time"

	pb "ride-sharing/shared/proto/user"
)

type UserModel struct {
	ID             string    `gorm:"column:id;type:uuid;primaryKey"`
	Username       string    `gorm:"column:username;type:varchar(255);not null"`
	ProfilePicture string    `gorm:"column:profile_picture;type:text"`
	Email          string    `gorm:"column:email;type:varchar(255);not null"`
	Password       string    `gorm:"column:password;type:varchar(255);not null"`
	Role           string    `gorm:"column:role;type:varchar(20);not null;default:rider"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

// ToProto converts a database User to a proto User
func (u *UserModel) ToProto() *pb.User {
	return &pb.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
}

// FromProtoUser converts a proto User to a database User
func FromProtoUser(u *pb.User) *UserModel {
	return &UserModel{
		ID:       u.Id,
		Username: u.Username,
		Email:    u.Email,
	}
}
