package in_memory_storage

import (
	"context"
	"time"
)

// Store is the in-memory key-value storage contract used by the domain.
// Callers depend only on this interface; the implementation may be Redis or anything else.
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}
