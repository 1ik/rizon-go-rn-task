package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"regexp"
)

var (
	// emailRegex is a simple regex pattern for basic email validation
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// generateEmailHash generates a SHA-256 hash from email and salt.
// The same email + salt will always produce the same hash (deterministic).
func generateEmailHash(email, salt string) string {
	// Concatenate email and salt
	data := email + salt

	// Compute SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Return hex-encoded hash string
	return hex.EncodeToString(hash[:])
}

// GenerateEmailAuthLink generates an email authentication link hash and stores it in Redis.
// It enforces rate limiting by checking if the email already exists in storage.
// Prints the link to the console and returns an error if something goes wrong.
func (a *appImpl) GenerateEmailAuthLink(ctx context.Context, email string) error {
	// Validate email format
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	// Check if email already exists in Redis (rate limiting)
	exists, err := a.store.Exists(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return ErrEmailAuthRateLimited
	}

	// Generate hash using SHA-256
	hash := generateEmailHash(email, a.authCfg.EmailAuthSalt)

	// Store email => hash in Redis with expiration
	if err := a.store.Set(ctx, email, hash, a.authCfg.EmailAuthExpiration); err != nil {
		return fmt.Errorf("failed to store email auth hash: %w", err)
	}

	// Build URI with endpoint: /{endpoint}?email={email}&secret={hash}
	values := url.Values{}
	values.Set("email", email)
	values.Set("secret", hash)
	uri := "/" + a.authCfg.EmailAuthEndpoint + "?" + values.Encode()

	// Print link to console
	log.Printf("Email auth link: %s", uri)

	return nil
}

// VerifyEmailAuth verifies an email authentication link by checking if the email exists
// in Redis and if the provided secret matches the stored hash.
func (a *appImpl) VerifyEmailAuth(ctx context.Context, email, secret string) error {
	// Validate email format
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	// Check if email exists in Redis
	exists, err := a.store.Exists(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if !exists {
		return ErrEmailAuthNotFound
	}

	// Get stored hash from Redis
	storedHash, err := a.store.Get(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get stored hash: %w", err)
	}

	// Compare provided secret with stored hash
	if secret != storedHash {
		return ErrEmailAuthInvalidSecret
	}

	// Verification successful - email remains in Redis (deletion handled by expiration)
	return nil
}

// isValidEmail performs basic email format validation.
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}
