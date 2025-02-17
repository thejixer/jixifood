package handlers

import (
	"context"
	"errors"

	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/logic"
	"github.com/thejixer/jixifood/services/auth/internal/utils"
	"github.com/thejixer/jixifood/shared/constants"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthLogicInterface interface {
	GenerateAndStoreOtp(phoneNumber string) (string, error)
	VerifyOTP(phoneNumber, otp string) (bool, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.UserEntity, error)
	CreateUser(ctx context.Context, phoneNumber, name string, roleID uint64) (*models.UserEntity, error)
	GetRequester(ctx context.Context) (*models.UserEntity, error)
	ConvertToPBUser(ctx context.Context, user *models.UserEntity) *pb.User
	CheckPermission(ctx context.Context, roleID uint64, permissionName string) bool
	ChangeUserRole(ctx context.Context, userID, roleID uint64) (*models.UserEntity, error)
	EditProfile(ctx context.Context, userID uint64, name string) (*models.UserEntity, error)
	QueryUsers(ctx context.Context, text string, page, limit uint64) ([]*pb.User, uint64, bool, error)
	GetUserByID(ctx context.Context, id uint64) (*models.UserEntity, error)
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

func (s *AuthHandler) VerifyOtp(ctx context.Context, req *pb.VerifyOtpRequest) (*pb.Token, error) {

	if len(req.Otp) < 4 {
		return nil, status.Error(codes.InvalidArgument, "bad request : otp wasn't properly provided")
	}

	normalizedPhone, err := utils.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+err.Error())
	}

	ok, err := s.AuthLogic.VerifyOTP(normalizedPhone, req.Otp)

	if err != nil || !ok {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, "bad request : "+err.Error())
		}
		if errors.Is(err, apperrors.ErrCodeMismatch) {
			return nil, status.Error(codes.Unauthenticated, "unauthorized : "+err.Error())
		}

		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())

	}

	user, err := s.AuthLogic.GetUserByPhoneNumber(ctx, normalizedPhone)
	if err != nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}

	var userId uint64
	if user == nil {
		// sign up logic
		user, err := s.AuthLogic.CreateUser(ctx, normalizedPhone, "", 0)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "bad request : "+err.Error())
		}
		userId = user.ID
	} else {
		// login logic
		userId = user.ID
	}

	tokenString, err := utils.SignToken(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}

	return &pb.Token{
		Token: tokenString,
	}, nil

}

func (s *AuthHandler) Me(ctx context.Context, req *pb.Empty) (*pb.User, error) {

	requester, err := s.AuthLogic.GetRequester(ctx)

	if err != nil {
		return nil, HandleGetRequesterError(err)
	}

	user := s.AuthLogic.ConvertToPBUser(ctx, requester)
	if user == nil {
		return nil, status.Error(codes.Internal, apperrors.ErrInternal.Error())
	}

	return user, nil

}

func (s *AuthHandler) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	requester, err := s.AuthLogic.GetRequester(ctx)

	if err != nil {
		return nil, HandleGetRequesterError(err)
	}

	ok := s.AuthLogic.CheckPermission(ctx, requester.RoleID, req.PersmissionName)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, apperrors.ErrForbidden.Error())
	}

	user := s.AuthLogic.ConvertToPBUser(ctx, requester)
	if user == nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}

	return &pb.CheckPermissionResponse{
		HasPermission: true,
		Requester:     user,
	}, nil
}

func (s *AuthHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {

	if req.RoleId == 0 || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+apperrors.ErrInputRequirements.Error())
	}
	normalizedPhone, err := utils.ValidatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+err.Error())
	}

	requester, err := s.AuthLogic.GetRequester(ctx)

	if err != nil {
		return nil, HandleGetRequesterError(err)
	}

	ok := s.AuthLogic.CheckPermission(ctx, requester.RoleID, constants.PermissionManageUser)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, apperrors.ErrForbidden.Error())
	}

	resp, err := s.AuthLogic.CreateUser(ctx, normalizedPhone, req.Name, req.RoleId)
	if err != nil {

		if errors.Is(err, apperrors.ErrDuplicatePhone) {
			return nil, status.Error(codes.InvalidArgument, apperrors.ErrDuplicatePhone.Error())
		}

		if errors.Is(err, apperrors.ErrInputRequirements) {
			return nil, status.Error(codes.InvalidArgument, apperrors.ErrInputRequirements.Error())
		}

		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}
	user := s.AuthLogic.ConvertToPBUser(ctx, resp)

	if user == nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}

	return user, nil
}

func (s *AuthHandler) ChangeUserRole(ctx context.Context, req *pb.ChangeUserRoleRequest) (*pb.User, error) {

	if req.UserId == 0 || req.RoleId == 0 {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+apperrors.ErrInputRequirements.Error())
	}

	requester, err := s.AuthLogic.GetRequester(ctx)

	if err != nil {
		return nil, HandleGetRequesterError(err)
	}

	ok := s.AuthLogic.CheckPermission(ctx, requester.RoleID, constants.PermissionManageUser)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, apperrors.ErrForbidden.Error())
	}

	resp, err := s.AuthLogic.ChangeUserRole(ctx, req.UserId, req.RoleId)
	if err != nil {
		if errors.Is(err, apperrors.ErrInputRequirements) || errors.Is(err, apperrors.ErrNotFound) {
			return nil, status.Error(codes.InvalidArgument, "bad request : "+apperrors.ErrInputRequirements.Error())
		}
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())

	}
	user := s.AuthLogic.ConvertToPBUser(ctx, resp)
	if user == nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}
	return user, nil
}

func (s *AuthHandler) EditProfile(ctx context.Context, req *pb.EditProfileRequest) (*pb.User, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+apperrors.ErrInputRequirements.Error())
	}
	requester, err := s.AuthLogic.GetRequester(ctx)

	if err != nil {
		return nil, HandleGetRequesterError(err)
	}

	resp, err := s.AuthLogic.EditProfile(ctx, requester.ID, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}
	user := s.AuthLogic.ConvertToPBUser(ctx, resp)
	if user == nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}

	return user, nil
}

func (s *AuthHandler) QueryUsers(ctx context.Context, req *pb.QueryUsersRequest) (*pb.QueryUsersResponse, error) {
	requester, err := s.AuthLogic.GetRequester(ctx)

	if err != nil {
		return nil, HandleGetRequesterError(err)
	}

	ok := s.AuthLogic.CheckPermission(ctx, requester.RoleID, constants.PermissionViewUser)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, apperrors.ErrForbidden.Error())
	}

	users, total, hasNextPage, err := s.AuthLogic.QueryUsers(ctx, req.Text, req.Page, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}

	return &pb.QueryUsersResponse{
		Data:        users,
		Total:       total,
		HasNextPage: hasNextPage,
	}, nil

}

func (s *AuthHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.User, error) {
	requester, err := s.AuthLogic.GetRequester(ctx)

	if err != nil {
		return nil, HandleGetRequesterError(err)
	}

	ok := s.AuthLogic.CheckPermission(ctx, requester.RoleID, constants.PermissionViewUser)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, apperrors.ErrForbidden.Error())
	}

	resp, err := s.AuthLogic.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, apperrors.ErrNotFound.Error())
	}
	user := s.AuthLogic.ConvertToPBUser(ctx, resp)
	if user == nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}
	return user, nil

}
