package handlers

import (
	"context"

	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/logic"
	"github.com/thejixer/jixifood/services/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthLogicInterface interface {
	GenerateAndStoreOtp(phoneNumber string) (string, error)
}

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	AuthLogic AuthLogicInterface
}

func NewAuthHandler(authLogic *logic.AuthLogic) *AuthHandler {
	return &AuthHandler{
		AuthLogic: authLogic,
	}
}

func (s *AuthHandler) RequestOtp(ctx context.Context, req *pb.RequestOtpRequest) (*pb.MessageResponse, error) {

	normalizedPhone, err := utils.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+err.Error())
	}

	otp, err := s.AuthLogic.GenerateAndStoreOtp(normalizedPhone)

	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate and store OTP")
	}

	smsError := utils.SendSMS(normalizedPhone, otp)
	if smsError != nil {
		return nil, status.Error(codes.Internal, "could not send the message to the customer")
	}

	return &pb.MessageResponse{
		Message: "ok",
	}, nil

}
