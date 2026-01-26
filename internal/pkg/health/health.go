package health

import (
	"net/http"
	"sync"
)

// State holds health status.
type State struct {
	mu      sync.RWMutex
	healthy bool
}

// NewState initializes health state.
func NewState() *State {
	return &State{healthy: true}
}

// SetHealthy updates the health state.
func (s *State) SetHealthy(value bool) {
	s.mu.Lock()
	s.healthy = value
	s.mu.Unlock()
}

// IsHealthy returns current health state.
func (s *State) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.healthy
}

// Handler builds HTTP handler for health endpoint.
func Handler(state *State) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if state.IsHealthy() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success\n"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failure\n"))
	}
}
