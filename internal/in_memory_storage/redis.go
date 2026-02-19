package in_memory_storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"rizon-test-task/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrNotConnected is returned when in-memory storage is used before NewStore().
	ErrNotConnected = errors.New("in-memory storage not connected")
)

var client *redis.Client

// connectRedis initializes the Redis connection using config. Idempotent.
func connectRedis() (*redis.Client, error) {
	if client != nil {
		return client, nil
	}

	cfg := config.GetRedisConfig()
	url := cfg.URL()

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client = redis.NewClient(opt)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Printf("In-memory storage (Redis) connected: %s", url)
	return client, nil
}

// redisStore implements Store using Redis.
type redisStore struct {
	client *redis.Client
}

func (s *redisStore) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

func (s *redisStore) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if err := s.client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (s *redisStore) Exists(ctx context.Context, key string) (bool, error) {
	count, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return count > 0, nil
}

func (s *redisStore) Delete(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// NewStore returns an in-memory Store, initializing the backend (Redis) under the hood.
// Callers use the returned Store without knowing which backend is used.
func NewStore() (Store, error) {
	c, err := connectRedis()
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotConnected
	}
	return &redisStore{client: c}, nil
}

// Close closes the in-memory storage connection. Call from shutdown.
func Close() error {
	if client == nil {
		return nil
	}
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close in-memory storage: %w", err)
	}
	client = nil
	return nil
}
