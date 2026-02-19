package app

import (
	"context"

	"rizon-test-task/internal/in_memory_storage"
	"rizon-test-task/internal/repository"
)

// App is the application business API. GraphQL and other adapters call only these methods.
type App interface {
	Hello(ctx context.Context) (string, error)
	Close() error
}

// appImpl holds wired dependencies and implements App.
type appImpl struct {
	userRepo repository.UserRepository
	store    in_memory_storage.Store
}

// New creates the app with provided dependencies.
// Returns the App interface. Call Close() when shutting down.
func New(userRepo repository.UserRepository, store in_memory_storage.Store) App {
	return &appImpl{
		userRepo: userRepo,
		store:    store,
	}
}

// Hello returns a greeting (business function).
func (a *appImpl) Hello(ctx context.Context) (string, error) {
	return "Hello, World!", nil
}

// Close releases all connections. Call from shutdown.
func (a *appImpl) Close() error {
	return in_memory_storage.Close()
}
