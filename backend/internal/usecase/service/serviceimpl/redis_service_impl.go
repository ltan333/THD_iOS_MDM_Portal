package serviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thienel/go-backend-template/internal/usecase/service"
	"github.com/thienel/go-backend-template/pkg/config"
)

type redisServiceImpl struct {
	client *redis.Client
}

// NewRedisService creates a new Redis service
func NewRedisService(cfg config.RedisConfig) service.RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &redisServiceImpl{
		client: client,
	}
}

func (s *redisServiceImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *redisServiceImpl) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *redisServiceImpl) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

func (s *redisServiceImpl) Exists(ctx context.Context, key string) (bool, error) {
	n, err := s.client.Exists(ctx, key).Result()
	return n > 0, err
}
