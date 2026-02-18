package config

import (
	"fmt"
	"os"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// GetServerConfig returns server configuration from environment or defaults
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: getEnv("PORT", "8080"),
	}
}

// Addr returns the listen address (e.g. ":8080")
func (c *ServerConfig) Addr() string {
	return ":" + c.Port
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetDatabaseConfig returns database configuration from environment or defaults
func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "rizon"),
		Password: getEnv("DB_PASSWORD", "rizon_dev_password"),
		DBName:   getEnv("DB_NAME", "rizon_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// DSN returns the PostgreSQL Data Source Name
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// URL returns the PostgreSQL connection URL
func (c *DatabaseConfig) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
