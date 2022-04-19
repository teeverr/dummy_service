package handlers

import (
	"log"
	"net/http"

	"github.com/teeverr/dummy_service/internal/health"
)

type ChecksControl struct {
	Liveness  *health.Probe
	Readiness *health.Probe
}

func (cc *ChecksControl) HealthOn(w http.ResponseWriter, _ *http.Request) {
	cc.Liveness.Enable()
	w.WriteHeader(http.StatusOK)
	log.Println("Health check endpoint enabled")
}

func (cc *ChecksControl) HealthOff(w http.ResponseWriter, _ *http.Request) {
	cc.Liveness.Disable()
	w.WriteHeader(http.StatusOK)
	log.Println("Health check endpoint disabled")
}

func (cc *ChecksControl) Ready(w http.ResponseWriter, _ *http.Request) {
	cc.Readiness.Disable()
	w.WriteHeader(http.StatusOK)
	log.Println("Readiness check endpoint enabled")
}

func (cc *ChecksControl) NotReady(w http.ResponseWriter, _ *http.Request) {
	cc.Readiness.Disable()
	w.WriteHeader(http.StatusOK)
	log.Println("Readiness check endpoint disabled")
}
