package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"rizon-test-task/internal/models"
	"rizon-test-task/internal/repository"
)

const (
	maxFeedbackContentLength = 5000 // Maximum length for feedback content
)

// SubmitFeedback creates a new feedback submission for the authenticated user.
// It validates the content, checks if feedback already exists for the user+device,
// and creates a new feedback record.
func (a *appImpl) SubmitFeedback(ctx context.Context, deviceID, content string) error {
	// Get authenticated user
	user, err := a.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	// Validate device ID
	if strings.TrimSpace(deviceID) == "" {
		return errors.New("device ID is required")
	}

	// Validate content
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("feedback content cannot be empty")
	}
	if len(content) > maxFeedbackContentLength {
		return fmt.Errorf("feedback content exceeds maximum length of %d characters", maxFeedbackContentLength)
	}

	// Check if feedback already exists for this user+device combination
	existingFeedback, err := a.feedbackRepo.GetByUserIDAndDeviceID(ctx, fmt.Sprintf("%d", user.ID), deviceID)
	if err != nil && !errors.Is(err, repository.ErrFeedbackNotFound) {
		return fmt.Errorf("failed to check existing feedback: %w", err)
	}
	if existingFeedback != nil {
		return ErrFeedbackAlreadySubmitted
	}

	// Create new feedback
	feedback := &models.Feedback{
		UserID:   user.ID,
		DeviceID: deviceID,
		Content:  content,
	}

	if err := a.feedbackRepo.Create(ctx, feedback); err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	return nil
}

// GetUserFeedbackOnDevice retrieves feedback for the authenticated user on a specific device.
// Returns the feedback if found, or nil, nil if not found (allows GraphQL query to return null).
func (a *appImpl) GetUserFeedbackOnDevice(ctx context.Context, deviceID string) (*models.Feedback, error) {
	// Get authenticated user
	user, err := a.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	// Validate device ID
	if strings.TrimSpace(deviceID) == "" {
		return nil, errors.New("device ID is required")
	}

	// Get feedback for user+device combination
	feedback, err := a.feedbackRepo.GetByUserIDAndDeviceID(ctx, fmt.Sprintf("%d", user.ID), deviceID)
	if err != nil {
		if errors.Is(err, repository.ErrFeedbackNotFound) {
			// Return nil, nil for not found (GraphQL query can return null)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return feedback, nil
}
