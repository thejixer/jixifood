package models

import (
	"context"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, phoneNumber, name string, roleId uint64) (*UserEntity, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*UserEntity, error)
	GetUserByID(ctx context.Context, id uint64) (*UserEntity, error)
	GetRoleByID(ctx context.Context, id uint64) (*Role, error)
	CheckPermission(ctx context.Context, roleID uint64, permissionName string) (bool, error)
	ChangeUserRole(ctx context.Context, userID, roleID uint64) (*UserEntity, error)
	EditProfile(ctx context.Context, userID uint64, name string) (*UserEntity, error)
}

type UserEntity struct {
	ID          uint64
	Name        string
	PhoneNumber string
	Status      string
	RoleID      uint64
	CreatedAt   time.Time
}

type User struct {
	ID          uint64
	Name        string
	PhoneNumber string
	Status      string
	Role        Role
	CreatedAt   time.Time
}

type Role struct {
	ID          uint64
	Name        string
	Description string
}

// ######### DTOS #########
type RequestOTPDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type VerifyOTPDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	OTP         string `json:"otp" validate:"required"`
}

type UserDto struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateUserDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Name        string `json:"name" validate:"required"`
	RoleID      uint64 `json:"roleID" validate:"required"`
}

type ChangeUserRoleDto struct {
	RoleID uint64 `json:"roleID" validate:"required"`
}

type EditProfileDto struct {
	Name string `json:"name" validate:"required"`
}
