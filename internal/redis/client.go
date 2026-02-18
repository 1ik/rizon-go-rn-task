package redis

import (
	"context"
	"fmt"
	"log"

	"rizon-test-task/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	// Client is the global Redis client
	Client *redis.Client
)

// Connect initializes the Redis connection
func Connect() (*redis.Client, error) {
	if Client != nil {
		return Client, nil
	}

	cfg := config.GetRedisConfig()
	url := cfg.URL()

	// Parse Redis URL
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Create Redis client
	Client = redis.NewClient(opt)

	// Test the connection
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Printf("Successfully connected to Redis: %s", url)
	return Client, nil
}

// Close closes the Redis connection
func Close() error {
	if Client == nil {
		return nil
	}

	if err := Client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	Client = nil
	return nil
}

// GetClient returns the global Redis client
func GetClient() *redis.Client {
	return Client
}
