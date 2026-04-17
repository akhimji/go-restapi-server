package observability

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// RequestIDKey is the context key for request ID
type RequestIDKey struct{}

// GenerateRequestID generates a new request ID as a hex string
func GenerateRequestID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// RequestIDMiddleware generates a request ID and adds it to the request context and response header
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new request ID
		requestID, err := GenerateRequestID()
		if err != nil {
			// If we can't generate an ID, we'll proceed without it but log the error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to request context
		ctx := context.WithValue(r.Context(), RequestIDKey{}, requestID)

		// Call the next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
