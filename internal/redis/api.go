package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrClientNotInitialized is returned when Redis operations are attempted
	// before calling Connect().
	ErrClientNotInitialized = errors.New("redis client not initialized")
)

// ensureClient checks if the Redis client is initialized
func ensureClient() error {
	if Client == nil {
		return ErrClientNotInitialized
	}
	return nil
}

// Get retrieves a string value from Redis
func Get(ctx context.Context, key string) (string, error) {
	if err := ensureClient(); err != nil {
		return "", err
	}

	val, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return val, nil
}

// Set stores a string value in Redis with optional expiration
func Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if err := ensureClient(); err != nil {
		return err
	}

	if err := Client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// Delete removes one or more keys from Redis
func Delete(ctx context.Context, keys ...string) error {
	if err := ensureClient(); err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	if err := Client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}

	return nil
}

// Exists checks if a key exists in Redis
func Exists(ctx context.Context, key string) (bool, error) {
	if err := ensureClient(); err != nil {
		return false, err
	}

	count, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}

	return count > 0, nil
}

// SetNX sets a key only if it does not already exist
func SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	if err := ensureClient(); err != nil {
		return false, err
	}

	set, err := Client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set key %s if not exists: %w", key, err)
	}

	return set, nil
}

// Incr increments a key's value by 1
func Incr(ctx context.Context, key string) (int64, error) {
	if err := ensureClient(); err != nil {
		return 0, err
	}

	val, err := Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}

	return val, nil
}

// Expire sets an expiration time on a key
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := ensureClient(); err != nil {
		return err
	}

	if err := Client.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration on key %s: %w", key, err)
	}

	return nil
}
