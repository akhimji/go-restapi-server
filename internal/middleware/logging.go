package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)
		next.ServeHTTP(w, req)
		log.Printf("Request completed in %v", time.Since(start))
	})
}
