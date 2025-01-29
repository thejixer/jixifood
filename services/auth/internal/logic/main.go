package logic

import (
	"context"
	"errors"

	"github.com/thejixer/jixifood/services/auth/internal/redis"
	"github.com/thejixer/jixifood/services/auth/internal/repository"
	"github.com/thejixer/jixifood/services/auth/internal/utils"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
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

func (l *AuthLogic) CreateUser(ctx context.Context, phoneNumber string, roleID uint64) (*models.UserEntity, error) {

	user, err := l.dbStore.AuthRepo.CreateUser(ctx, phoneNumber, roleID)
	if err != nil {
		return nil, err
	}
	return user, nil

}
