package handlers

import (
	"log"
	"net/http"
)

type HealthCheck struct {
	l *log.Logger
}

// NewHealthCheck returns a new HealthCheck
func NewHealthCheck(l *log.Logger) *HealthCheck {
	hh := new(HealthCheck)
	hh.l = l
	return hh
}

// swagger:route GET /health heath checkHealthStatus
// Returns no content and checks the health status
// responses:
//	200: noContentResponse

// CheckHealthStatus checks the health status
func (h *HealthCheck) CheckHealthStatus(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("[DEBUG] Get check health status")

	rw.WriteHeader(http.StatusOK)
}
