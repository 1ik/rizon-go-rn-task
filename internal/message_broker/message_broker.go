package message_broker

import (
	"context"
)

// MessageBroker is the message broker contract used by the domain.
// Callers depend only on this interface; the implementation may be RabbitMQ or anything else.
type MessageBroker interface {
	PublishEmailJob(ctx context.Context, to, subject, body string) error
	Close() error
}
