package app

import (
	"context"

	"rizon-test-task/internal/database"
	"rizon-test-task/internal/in_memory_storage"

	"gorm.io/gorm"
)

// App is the application business API. GraphQL and other adapters call only these methods.
type App interface {
	Hello(ctx context.Context) (string, error)
	Close() error
}

// appImpl holds wired dependencies and implements App.
type appImpl struct {
	db    *gorm.DB
	store in_memory_storage.Store
}

// New creates and wires all dependencies (database, in-memory storage backed by Redis).
// Returns the App interface. Call Close() when shutting down.
func New() (App, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	store, err := in_memory_storage.NewStore()
	if err != nil {
		_ = database.Close()
		return nil, err
	}

	return &appImpl{
		db:    db,
		store: store,
	}, nil
}

// Hello returns a greeting (business function).
func (a *appImpl) Hello(ctx context.Context) (string, error) {
	return "Hello, World!", nil
}

// Close releases all connections. Call from shutdown.
func (a *appImpl) Close() error {
	if err := in_memory_storage.Close(); err != nil {
		return err
	}
	return database.Close()
}
