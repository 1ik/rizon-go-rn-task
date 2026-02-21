package postgres

import (
	"context"
	"fmt"

	"rizon-test-task/internal/models"
	"rizon-test-task/internal/repository"

	"gorm.io/gorm"
)

// feedbackRepository implements repository.FeedbackRepository using PostgreSQL.
type feedbackRepository struct {
	db *gorm.DB
}

// NewFeedbackRepository returns a PostgreSQL implementation of FeedbackRepository.
func NewFeedbackRepository(db *gorm.DB) repository.FeedbackRepository {
	return &feedbackRepository{db: db}
}

func (r *feedbackRepository) Create(ctx context.Context, feedback *models.Feedback) error {
	if err := r.db.WithContext(ctx).Create(feedback).Error; err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}
	return nil
}
