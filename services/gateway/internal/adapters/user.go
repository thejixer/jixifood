package adapters

import (
	authPB "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/shared/models"
)

func MapPBUserToUserDTO(user *authPB.User) *models.UserDto {

	return &models.UserDto{
		ID:          user.Id,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Status:      user.Status.String(),
		Role:        user.Role,
		CreatedAt:   user.CreatedAt.AsTime(),
	}
}
