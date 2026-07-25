package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Haarish-Ahmad/idempotent-gateway/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		Password: cfg.RedisPasseord,
		DB: cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	fmt.Printf("Redis Connection Established: %s\n", pong)
	return &RedisClient{Client: rdb}, nil
}