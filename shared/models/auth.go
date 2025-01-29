package models

import (
	"context"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, phoneNumber string, roleId uint64) (*UserEntity, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*UserEntity, error)
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
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ######### DTOS #########
type RequestOTPDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type VerifyOTPDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	OTP         string `json:"otp" validate:"required"`
}
