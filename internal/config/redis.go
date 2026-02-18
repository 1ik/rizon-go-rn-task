package config

import (
	"fmt"
	"os"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URLString string
	Host      string
	Port      string
	Password  string
	DB        int
}

// GetRedisConfig returns Redis configuration from environment or defaults
func GetRedisConfig() *RedisConfig {
	// Check if REDIS_URL is set, use it directly
	if url := os.Getenv("REDIS_URL"); url != "" {
		return &RedisConfig{
			URLString: url,
		}
	}

	// Otherwise, build from individual components
	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0, // Default DB 0
	}
}

// URL returns the Redis connection URL
func (c *RedisConfig) URL() string {
	if c.URLString != "" {
		return c.URLString
	}

	if c.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%d", c.Password, c.Host, c.Port, c.DB)
	}
	return fmt.Sprintf("redis://%s:%s/%d", c.Host, c.Port, c.DB)
}
