package config

import (
	"time"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	EmailAuthSalt       string
	EmailAuthExpiration time.Duration
	EmailAuthEndpoint   string
}

// GetAuthConfig returns authentication configuration from environment or defaults
func GetAuthConfig() *AuthConfig {
	salt := getEnv("EMAIL_AUTH_SALT", "")
	if salt == "" {
		// In production, this should be required, but for development we provide a default
		// TODO: Consider making this required in production
		salt = "default-dev-salt-change-in-production"
	}

	endpoint := getEnv("EMAIL_AUTH_ENDPOINT", "email-auth")

	return &AuthConfig{
		EmailAuthSalt:       salt,
		EmailAuthExpiration: 30 * time.Minute,
		EmailAuthEndpoint:   endpoint,
	}
}
