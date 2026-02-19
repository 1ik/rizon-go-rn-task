package graphql

import (
	"rizon-test-task/internal/app"
)

// Resolver is the GraphQL resolver. It delegates to the app (API gateway only).
type Resolver struct {
	App app.App
}

// NewResolver returns a resolver that uses the given app for all business operations.
func NewResolver(a app.App) *Resolver {
	return &Resolver{App: a}
}
