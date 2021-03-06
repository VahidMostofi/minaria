package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type HealthCheck struct {
	l *log.Logger
}

func (h *HealthCheck) AttachRouter(mr *mux.Router) *mux.Router {
	heathHandler := mr.PathPrefix("/health").Subrouter()
	heathHandler.HandleFunc("", h.CheckHealthStatus).Methods(http.MethodGet)
	return heathHandler
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
