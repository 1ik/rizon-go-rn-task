package app

import (
	"context"
	"errors"

	"rizon-test-task/internal/config"
	"rizon-test-task/internal/in_memory_storage"
	"rizon-test-task/internal/message_broker"
	"rizon-test-task/internal/repository"
)

var (
	// ErrEmailAuthRateLimited is returned when an email auth link has already been requested
	// and the rate limit period has not expired yet.
	ErrEmailAuthRateLimited = errors.New("email auth link already requested, please try again later")
	// ErrEmailAuthNotFound is returned when the email auth link has expired or doesn't exist.
	ErrEmailAuthNotFound = errors.New("email auth link not found or expired")
	// ErrEmailAuthInvalidSecret is returned when the secret doesn't match the stored hash.
	ErrEmailAuthInvalidSecret = errors.New("invalid secret")
)

// App is the application business API. GraphQL and other adapters call only these methods.
type App interface {
	Hello(ctx context.Context) (string, error)
	GenerateEmailAuthLink(ctx context.Context, email string) error
	VerifyEmailAuth(ctx context.Context, email, secret string) (string, error)
	Close() error
}

// appImpl holds wired dependencies and implements App.
type appImpl struct {
	userRepo      repository.UserRepository
	store         in_memory_storage.Store
	authCfg       *config.AuthConfig
	messageBroker message_broker.MessageBroker
}

// New creates the app with provided dependencies.
// Returns the App interface. Call Close() when shutting down.
func New(userRepo repository.UserRepository, store in_memory_storage.Store, authCfg *config.AuthConfig, messageBroker message_broker.MessageBroker) App {
	return &appImpl{
		userRepo:      userRepo,
		store:         store,
		authCfg:       authCfg,
		messageBroker: messageBroker,
	}
}

// Hello returns a greeting (business function).
func (a *appImpl) Hello(ctx context.Context) (string, error) {
	return "Hello, World!", nil
}

// Close releases all connections. Call from shutdown.
func (a *appImpl) Close() error {
	if err := a.messageBroker.Close(); err != nil {
		return err
	}
	return in_memory_storage.Close()
}
