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


// FeedbackRepository defines operations for feedback data access.
type FeedbackRepository interface {
	Create(ctx context.Context, feedback *models.Feedback) error
}
