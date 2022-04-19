package health

import (
	"net/http"
	"sync/atomic"
)

type Probe struct {
	name   string
	status int32
}

func (p *Probe) Status() bool {
	return atomic.LoadInt32(&p.status) == 1
}

func (p *Probe) Enable() {
	atomic.StoreInt32(&p.status, 1)
}

func (p *Probe) Disable() {
	atomic.StoreInt32(&p.status, 0)
}

func (p *Probe) Name() string {
	return p.name
}

func NewProbe(name string, status int32) *Probe {
	return &Probe{name, status}
}

func (h *Probe) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if h.Status() {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))

		return
	}

	w.WriteHeader(http.StatusServiceUnavailable)
}

func NewLivenessProbe() *Probe {
	return NewProbe("liveness", 1)
}

func NewReadinessProbe() *Probe {
	return NewProbe("readiness", 0)
}
