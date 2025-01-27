package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/thejixer/jixifood/services/auth/internal/config"
)

type RedisStore struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisStore(c *config.AuthServiceConfig) (*RedisStore, error) {
	var ctx = context.Background()

	Addr := c.RedisURI

	rdb := redis.NewClient(&redis.Options{
		Addr: Addr,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		rdb: rdb,
		ctx: ctx,
	}, nil
}
