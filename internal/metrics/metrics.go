package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics holds the runtime metrics
type Metrics struct {
	startTime           time.Time
	totalRequests       uint64
	requestsByStatus    map[string]uint64
	requestsByStatusMu  sync.RWMutex
	peopleCount         int64
	peopleCountMu       sync.RWMutex
}

// NewMetrics creates and returns a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		startTime:        time.Now(),
		requestsByStatus: make(map[string]uint64),
	}
}

// GetStartTime returns the server start time
func (m *Metrics) GetStartTime() time.Time {
	return m.startTime
}

// IncrementTotalRequests increments the total requests counter
func (m *Metrics) IncrementTotalRequests() {
	atomic.AddUint64(&m.totalRequests, 1)
}

// IncrementRequestsByStatus increments the counter for a specific status code
func (m *Metrics) IncrementRequestsByStatus(statusCode string) {
	m.requestsByStatusMu.Lock()
	defer m.requestsByStatusMu.Unlock()
	m.requestsByStatus[statusCode]++
}

// GetTotalRequests returns the total number of requests handled
func (m *Metrics) GetTotalRequests() uint64 {
	return atomic.LoadUint64(&m.totalRequests)
}

// GetRequestsByStatus returns a copy of the requests by status map
func (m *Metrics) GetRequestsByStatus() map[string]uint64 {
	m.requestsByStatusMu.RLock()
	defer m.requestsByStatusMu.RUnlock()

	// Return a copy to avoid external mutation
	result := make(map[string]uint64)
	for k, v := range m.requestsByStatus {
		result[k] = v
	}
	return result
}

// SetPeopleCount sets the current people count
func (m *Metrics) SetPeopleCount(count int64) {
	m.peopleCountMu.Lock()
	defer m.peopleCountMu.Unlock()
	m.peopleCount = count
}

// GetPeopleCount returns the current people count
func (m *Metrics) GetPeopleCount() int64 {
	m.peopleCountMu.RLock()
	defer m.peopleCountMu.RUnlock()
	return m.peopleCount
}