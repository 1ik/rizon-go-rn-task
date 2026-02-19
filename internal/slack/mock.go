package slack

import (
	"context"
	"log"
)

// mockSlackSender implements SlackSender by logging messages to console.
type mockSlackSender struct{}

// NewMockSender returns a mock implementation of SlackSender that logs to console.
func NewMockSender() SlackSender {
	return &mockSlackSender{}
}

// SendMessage logs the Slack message to console.
func (m *mockSlackSender) SendMessage(ctx context.Context, message string) error {
	log.Printf("[MOCK SLACK] Message:\n%s", message)
	return nil
}
