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
// Publishes the email job to the message broker after successful link generation.
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

	// Build full URL with endpoint: {baseURL}/{endpoint}?email={email}&secret={hash}
	values := url.Values{}
	values.Set("email", email)
	values.Set("secret", hash)
	path := "/" + a.authCfg.EmailAuthEndpoint + "?" + values.Encode()
	uri := a.authCfg.BaseURL + path

	// Publish email job to message broker
	if err := a.publishEmailAuthJob(ctx, email, uri); err != nil {
		// If publishing fails, clean up the stored email entry so user can retry
		if delErr := a.store.Delete(ctx, email); delErr != nil {
			// Log cleanup error but return the original publishing error
			log.Printf("Warning: failed to cleanup email entry after publish failure: %v", delErr)
		}
		return fmt.Errorf("failed to publish email job: %w", err)
	}

	return nil
}

// VerifyEmailAuth verifies an email authentication link by checking if the email exists
// in Redis and if the provided secret matches the stored hash.
// Returns the auth link URL on successful verification.
func (a *appImpl) VerifyEmailAuth(ctx context.Context, email, secret string) (string, error) {
	// Validate email format
	if !isValidEmail(email) {
		return "", fmt.Errorf("invalid email format")
	}

	// Check if email exists in Redis
	exists, err := a.store.Exists(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to check email existence: %w", err)
	}
	if !exists {
		return "", ErrEmailAuthNotFound
	}

	// Get stored hash from Redis
	storedHash, err := a.store.Get(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get stored hash: %w", err)
	}

	// Compare provided secret with stored hash
	if secret != storedHash {
		return "", ErrEmailAuthInvalidSecret
	}

	// Build and return the auth link URL
	values := url.Values{}
	values.Set("email", email)
	values.Set("secret", secret)
	path := "/" + a.authCfg.EmailAuthEndpoint + "?" + values.Encode()
	uri := a.authCfg.BaseURL + path

	// Verification successful - email remains in Redis (deletion handled by expiration)
	return uri, nil
}

// publishEmailAuthJob publishes an email authentication job to the message broker.
// Returns an error if publishing fails.
func (a *appImpl) publishEmailAuthJob(ctx context.Context, email, authLink string) error {
	subject := "Rizon: Your App Authentication Link"
	body := fmt.Sprintf(`
		<p>Please Click the link below to sign in to your app:</p>
		<p><a href="%s">%s</a></p>
		<p>This link will expire in 30 minutes.</p>
	`, authLink, authLink)

	if err := a.messageBroker.PublishEmailJob(ctx, email, subject, body); err != nil {
		return fmt.Errorf("failed to publish email job: %w", err)
	}
	return nil
}

// isValidEmail performs basic email format validation.
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}
