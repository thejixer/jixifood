package handlers

import (
	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/redis"
	"github.com/thejixer/jixifood/services/auth/internal/repository"
)

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
	dbStore    *repository.PostgresStore
	RedisStore *redis.RedisStore
}

func NewAuthServiceServer(dbStore *repository.PostgresStore, redisStore *redis.RedisStore) *authServiceServer {
	return &authServiceServer{
		dbStore:    dbStore,
		RedisStore: redisStore,
	}
}
