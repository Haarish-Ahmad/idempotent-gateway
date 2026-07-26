package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Haarish-Ahmad/idempotent-gateway/internal/config"
	"github.com/Haarish-Ahmad/idempotent-gateway/internal/models"
	"github.com/redis/go-redis/v9"
)

//----------------------------------------------------------------------------------------------------------------
var (
	ErrLockConflict  = errors.New("request is currently in progress")
	ErrAlreadyCached = errors.New("request already completed and cached")
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Redis Connection Established: PONG")
	return &RedisClient{Client: rdb}, nil
}

//----------------------------------------------------------------------------------------------------------------


func (r *RedisClient) EvalIdempotencyLock(ctx context.Context, key string, ttl time.Duration) (*models.LockResult, error) {
	acquired, err := r.Client.SetNX(ctx, key, string(models.StateInProgress), ttl).Result()
	if err != nil {
		return nil, fmt.Errorf("redis setnx failed: %w", err)
	}
	if acquired {
		return &models.LockResult{State: models.StateNotFound}, nil
	}

	existing, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis get failed after setnx: %w", err)
	}

	if existing == string(models.StateInProgress) {
		return &models.LockResult{State: models.StateInProgress}, ErrLockConflict
	}

	return &models.LockResult{
		State: models.StateCompleted,
		CachedResponse: &models.CachedResponse{
			Body: []byte(existing),
		},
	}, ErrAlreadyCached
}

func (r *RedisClient) SaveResponse(ctx context.Context, key string, payload string, ttl time.Duration) error {
	return r.Client.Set(ctx, key, payload, ttl).Err()
}