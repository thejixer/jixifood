package logic

import (
	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/shared/constants"
)

func GetUserStatus(status string) pb.UserStatus {
	switch status {
	case constants.UserStatusComplete:
		return pb.UserStatus_complete
	case constants.UserStatusIncomplete:
		return pb.UserStatus_incomplete
	default:
		return pb.UserStatus_incomplete
	}
}
