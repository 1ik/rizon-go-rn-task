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
	JWTSecret           string
	JWTExpiration       time.Duration
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

	// JWT configuration
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		// In production, this should be required
		jwtSecret = "default-dev-jwt-secret-change-in-production"
	}

	return &AuthConfig{
		EmailAuthSalt:       salt,
		EmailAuthExpiration: 5 * time.Minute,
		EmailAuthEndpoint:   endpoint,
		BaseURL:             baseURL,
		JWTSecret:           jwtSecret,
		JWTExpiration:       7 * 24 * time.Hour, // 7 days
	}
}
