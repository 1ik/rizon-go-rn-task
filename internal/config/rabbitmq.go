package config

import (
	"fmt"
	"os"
)

// RabbitMQConfig holds RabbitMQ configuration
type RabbitMQConfig struct {
	URLString string
	Host      string
	Port      string
	User      string
	Password  string
	VHost     string
}

// GetRabbitMQConfig returns RabbitMQ configuration from environment or defaults
func GetRabbitMQConfig() *RabbitMQConfig {
	// Check if RABBITMQ_URL is set, use it directly
	if url := os.Getenv("RABBITMQ_URL"); url != "" {
		return &RabbitMQConfig{
			URLString: url,
		}
	}

	// Otherwise, build from individual components
	return &RabbitMQConfig{
		Host:     getEnv("RABBITMQ_HOST", "localhost"),
		Port:     getEnv("RABBITMQ_PORT", "5672"),
		User:     getEnv("RABBITMQ_USER", "rizon"),
		Password: getEnv("RABBITMQ_PASSWORD", "rizon_dev_password"),
		VHost:    getEnv("RABBITMQ_VHOST", "/"),
	}
}

// URL returns the RabbitMQ connection URL
func (c *RabbitMQConfig) URL() string {
	if c.URLString != "" {
		return c.URLString
	}

	return fmt.Sprintf("amqp://%s:%s@%s:%s%s", c.User, c.Password, c.Host, c.Port, c.VHost)
}
