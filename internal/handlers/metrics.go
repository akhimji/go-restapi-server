package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-restapi-server/internal/metrics"
	"go-restapi-server/internal/store"
)

// MetricsHandler handles metrics related endpoints
type MetricsHandler struct {
	metricsStore *metrics.Metrics
	store        store.PersonStore
}

// NewMetricsHandler creates a new MetricsHandler
func NewMetricsHandler(metricsStore *metrics.Metrics, store store.PersonStore) *MetricsHandler {
	return &MetricsHandler{
		metricsStore: metricsStore,
		store:        store,
	}
}

// MetricsEndpoint returns runtime metrics as JSON
func (h *MetricsHandler) MetricsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Get current people count
	peopleCount := int64(len(h.store.List()))

	// Update metrics with current count
	h.metricsStore.SetPeopleCount(peopleCount)

	// Return metrics as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"uptime_seconds":   time.Since(h.metricsStore.GetStartTime()).Seconds(),
		"total_requests":   h.metricsStore.GetTotalRequests(),
		"requests_by_status": h.metricsStore.GetRequestsByStatus(),
		"people_count":     peopleCount,
	})
}