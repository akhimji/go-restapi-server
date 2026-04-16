package observability

import (
	"log/slog"
	"net/http"
	"time"
)

// ResponseWriterWrapper wraps http.ResponseWriter to capture status code
type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (w *ResponseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// NewResponseWriterWrapper creates a new ResponseWriterWrapper
func NewResponseWriterWrapper(w http.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter: w,
		statusCode:     200, // Default status code
	}
}

// LoggingMiddleware logs requests with structured JSON output
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap the response writer to capture status code
		wrapper := NewResponseWriterWrapper(w)
		
		// Get request ID from context, if available
		var requestID string
		if reqID, ok := r.Context().Value(RequestIDKey{}).(string); ok {
			requestID = reqID
		}
		
		// Call the next handler
		next.ServeHTTP(wrapper, r)
		
		// Calculate duration
		duration := time.Since(start)
		
		// Log the request
		slog.Info("Request completed",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrapper.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}