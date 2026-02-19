package config

import (
	"fmt"
	"log"
	"strconv"
)

// EmailConfig holds email/SMTP configuration
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	FromEmail string
}

// GetEmailConfig returns email configuration from environment or defaults
func GetEmailConfig() *EmailConfig {
	host := getEnv("EMAIL_SMTP_HOST", "smtp.gmail.com")
	
	portStr := getEnv("EMAIL_SMTP_PORT", "587")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		// Default to 587 if invalid port
		port = 587
	}

	username := getEnv("EMAIL_USERNAME", "")
	if username == "" {
		log.Fatal("EMAIL_USERNAME environment variable is required")
	}

	password := getEnv("EMAIL_PASSWORD", "")
	if password == "" {
		log.Fatal("EMAIL_PASSWORD environment variable is required")
	}

	fromEmail := getEnv("EMAIL_FROM", "")
	if fromEmail == "" {
		log.Fatal("EMAIL_FROM environment variable is required")
	}

	return &EmailConfig{
		SMTPHost: host,
		SMTPPort: port,
		Username: username,
		Password: password,
		FromEmail: fromEmail,
	}
}

// Addr returns the SMTP server address (e.g., "smtp.gmail.com:587")
func (c *EmailConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.SMTPHost, c.SMTPPort)
}
