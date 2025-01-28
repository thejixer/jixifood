package models

type AuthRepository interface{}

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ######### DTOS #########
type RequestOTPDto struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}
