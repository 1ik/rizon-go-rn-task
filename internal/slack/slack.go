package slack

import (
	"context"
)

// SlackSender is the Slack sending contract used by the domain.
// Callers depend only on this interface; the implementation may be HTTP webhook or anything else.
type SlackSender interface {
	SendMessage(ctx context.Context, message string) error
}
