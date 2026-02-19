package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"rizon-test-task/internal/models"
	"rizon-test-task/internal/repository"

	"github.com/golang-jwt/jwt/v5"
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

// EmailAuth verifies an email authentication link by checking if the email exists
// in Redis and if the provided secret matches the stored hash.
// If verification succeeds, it creates the user if they don't exist and returns a JWT token.
func (a *appImpl) EmailAuth(ctx context.Context, email, secret string) (string, error) {
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

	// Verification successful - find or create user
	user, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			// User doesn't exist, create them
			user = &models.User{
				Email: email,
			}
			if err := a.userRepo.Create(ctx, user); err != nil {
				return "", fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to find user: %w", err)
		}
	}

	// Generate JWT token
	token, err := a.generateJWT(user.ID, user.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return token, nil
}

// generateJWT creates a JWT token with user ID and email claims.
func (a *appImpl) generateJWT(userID uint, email string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": strconv.FormatUint(uint64(userID), 10),
		"email":   email,
		"iat":     now.Unix(),
		"exp":     now.Add(a.authCfg.JWTExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.authCfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and extracts user ID and email from claims.
// Returns user ID, email, and error.
func (a *appImpl) ValidateJWT(ctx context.Context, tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.authCfg.JWTSecret), nil
	})

	if err != nil {
		return 0, "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, "", ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	// Extract user_id
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return 0, "", fmt.Errorf("user_id claim not found or invalid")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, "", fmt.Errorf("invalid user_id format: %w", err)
	}

	// Extract email
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return 0, "", fmt.Errorf("email claim not found or invalid")
	}

	return uint(userID), email, nil
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

// GetCurrentUser retrieves the current authenticated user from the context.
// The context should contain a JWT token that has been validated by middleware.
func (a *appImpl) GetCurrentUser(ctx context.Context) (*models.User, error) {
	// Extract token from context (set by middleware)
	// Use the same key as middleware: "token"
	tokenString, ok := ctx.Value("token").(string)
	if !ok || tokenString == "" {
		return nil, ErrUnauthorized
	}

	// Validate JWT and extract user ID
	userID, _, err := a.ValidateJWT(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	// Fetch user from repository
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	user, err := a.userRepo.FindByID(ctx, userIDStr)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// isValidEmail performs basic email format validation.
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}
