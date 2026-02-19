package graphql

import (
	"context"
	"net/http"
	"strings"
)

// tokenKey is the context key for storing JWT token
// Using a string key allows it to be accessed from other packages
const tokenKey = "token"

// AuthMiddleware extracts JWT token from Authorization header and adds it to context.
// It does not require authentication for all requests - public queries can proceed without token.
// The token is extracted if present and validated by resolvers that need it.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		
		// If Authorization header is present, extract Bearer token
		if authHeader != "" {
			// Check if it's a Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token := parts[1]
				// Add token to context using string key
				ctx := context.WithValue(r.Context(), tokenKey, token)
				r = r.WithContext(ctx)
			}
		}
		
		// Continue to next handler (token extraction is optional)
		next.ServeHTTP(w, r)
	})
}

// GetTokenFromContext extracts the JWT token from context.
func GetTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenKey).(string)
	return token, ok && token != ""
}
