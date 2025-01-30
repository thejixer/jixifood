package logic

import (
	"context"
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt"
	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/redis"
	"github.com/thejixer/jixifood/services/auth/internal/repository"
	"github.com/thejixer/jixifood/services/auth/internal/utils"
	"github.com/thejixer/jixifood/shared/constants"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthLogic struct {
	dbStore    *repository.PostgresStore
	RedisStore *redis.RedisStore
}

// this is the service layer, but since the term service is repeated in different meanings,
// I renamed it to AuthLogic instead of AuthService for clarity
func NewAuthLogic(dbStore *repository.PostgresStore, redisStore *redis.RedisStore) *AuthLogic {
	return &AuthLogic{
		dbStore:    dbStore,
		RedisStore: redisStore,
	}
}

func (l *AuthLogic) GenerateAndStoreOtp(phoneNumber string) (string, error) {

	otp := utils.GenerateOTP()
	err := l.RedisStore.SetOTP(phoneNumber, otp)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (l *AuthLogic) VerifyOTP(phoneNumber, otp string) (bool, error) {
	generatedOTP, err := l.RedisStore.GetOTP(phoneNumber)
	if err != nil {
		return false, apperrors.ErrNotFound
	}

	if generatedOTP != otp {
		return false, apperrors.ErrCodeMismatch
	}

	go l.RedisStore.DelOTP(phoneNumber)

	return true, nil
}

func (l *AuthLogic) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.UserEntity, error) {

	user, err := l.dbStore.AuthRepo.GetUserByPhoneNumber(ctx, phoneNumber)

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			// this behaviour is expected
			return nil, nil
		}
		return nil, apperrors.ErrUnexpected
	}

	return user, nil
}

func (l *AuthLogic) CreateUser(ctx context.Context, phoneNumber, name string, roleID uint64) (*models.UserEntity, error) {

	user, err := l.dbStore.AuthRepo.CreateUser(ctx, phoneNumber, name, roleID)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (l *AuthLogic) GetRequester(ctx context.Context) (*models.UserEntity, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, apperrors.ErrMissingMetaData
	}
	tokens := md["auth"]
	if len(tokens) == 0 {
		return nil, apperrors.ErrMissingToken
	}
	tokenString := tokens[0]
	token, err := utils.VerifyToken(tokenString)

	if err != nil || !token.Valid {
		return nil, apperrors.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["id"] == nil {
		return nil, apperrors.ErrUnauthorized
	}
	i := claims["id"].(string)
	intId, err := strconv.Atoi(i)
	if err != nil {
		return nil, apperrors.ErrInternal
	}
	userId := uint64(intId)

	userFromCache := l.RedisStore.GetUser(userId)
	if userFromCache != nil {
		return userFromCache, nil
	}

	user, err := l.dbStore.AuthRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	go l.RedisStore.CacheUser(user)

	return user, nil

}

func (l *AuthLogic) ConvertToPBUser(ctx context.Context, user *models.UserEntity) *pb.User {

	role, err := l.dbStore.AuthRepo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		return nil
	}

	userStatus := func(status string) pb.UserStatus {
		switch status {
		case constants.UserStatusComplete:
			return pb.UserStatus_complete
		case constants.UserStatusIncomplete:
			return pb.UserStatus_incomplete
		default:
			return pb.UserStatus_incomplete
		}
	}(user.Status)

	return &pb.User{
		Id:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Status:      userStatus,
		Role:        role.Name,
		CreatedAt:   timestamppb.New(user.CreatedAt),
	}
}

func (l *AuthLogic) CheckPermission(ctx context.Context, roleID uint64, permissionName string) bool {

	ok, err := l.dbStore.AuthRepo.CheckPermission(ctx, roleID, permissionName)
	if err != nil || !ok {
		return false
	}

	return true
}
