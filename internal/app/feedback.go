package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"rizon-test-task/internal/models"
)

const (
	maxFeedbackContentLength = 5000 // Maximum length for feedback content
)

// SubmitFeedback creates a new feedback submission for the authenticated user.
func (a *appImpl) SubmitFeedback(ctx context.Context, deviceID, content string) error {
	user, err := a.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	if strings.TrimSpace(deviceID) == "" {
		return errors.New("device ID is required")
	}

	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("feedback content cannot be empty")
	}
	if len(content) > maxFeedbackContentLength {
		return fmt.Errorf("feedback content exceeds maximum length of %d characters", maxFeedbackContentLength)
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

	// Publish Slack job after successful feedback creation
	// Handle errors gracefully - log but don't fail feedback creation
	if err := a.publishSlackJob(ctx, user.Email, deviceID, content, feedback.ID); err != nil {
		log.Printf("Warning: failed to publish slack job for feedback ID %d: %v", feedback.ID, err)
		// Don't return error - feedback was successfully created
	}

	return nil
}

// publishSlackJob publishes a Slack notification job to the message broker.
// Returns an error if publishing fails.
func (a *appImpl) publishSlackJob(ctx context.Context, userEmail, deviceID, content string, feedbackID uint) error {
	if err := a.messageBroker.PublishSlackJob(ctx, userEmail, deviceID, content, feedbackID); err != nil {
		return fmt.Errorf("failed to publish slack job: %w", err)
	}
	return nil
}
