package logic

import (
	"github.com/thejixer/jixifood/services/auth/internal/redis"
	"github.com/thejixer/jixifood/services/auth/internal/repository"
	"github.com/thejixer/jixifood/services/auth/internal/utils"
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
