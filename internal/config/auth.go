package config

import (
	"time"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	EmailAuthSalt       string
	EmailAuthExpiration time.Duration
	EmailAuthEndpoint   string
	BaseURL             string
}

// GetAuthConfig returns authentication configuration from environment or defaults
func GetAuthConfig(serverCfg *ServerConfig) *AuthConfig {
	salt := getEnv("EMAIL_AUTH_SALT", "")
	if salt == "" {
		// In production, this should be required, but for development we provide a default
		// TODO: Consider making this required in production
		salt = "default-dev-salt-change-in-production"
	}

	endpoint := getEnv("EMAIL_AUTH_ENDPOINT", "email-auth")

	// Use BaseURL from server config (which handles env var or defaults to localhost:port)
	baseURL := serverCfg.BaseURL

	return &AuthConfig{
		EmailAuthSalt:       salt,
		EmailAuthExpiration: 1 * time.Minute,
		EmailAuthEndpoint:   endpoint,
		BaseURL:             baseURL,
	}
}
