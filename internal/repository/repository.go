package repository

import (
	"context"
	"errors"
	"rizon-test-task/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository defines operations for user data access.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

// LoginRepository defines operations for login/authentication data access.
type LoginRepository interface {
	// Add login-related methods here as needed
	// Example: StoreMagicLink(ctx, email, hash, expiresAt) error
}
