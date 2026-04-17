package handlers

import (
	"net/http"
)

// HealthEndpoint returns the health status of the service
func HealthEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
