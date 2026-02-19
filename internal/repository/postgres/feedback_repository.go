package postgres

import (
	"context"
	"errors"
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

func (r *feedbackRepository) GetByUserIDAndDeviceID(ctx context.Context, userID string, deviceID string) (*models.Feedback, error) {
	var feedback models.Feedback
	if err := r.db.WithContext(ctx).Where("user_id = ? AND device_id = ?", userID, deviceID).First(&feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrFeedbackNotFound
		}
		return nil, fmt.Errorf("failed to get feedback by user id and device id: %w", err)
	}
	return &feedback, nil
}
