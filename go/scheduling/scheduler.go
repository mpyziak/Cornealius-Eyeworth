package scheduling

import (
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

// Scheduler wraps a robfig/cron instance and supports hot-restart when the
// CRON expression changes.
type Scheduler struct {
	mu      sync.Mutex
	crontab *cron.Cron
}

// Start stops any existing schedule and begins a new one. onFire is called each
// time the expression triggers. The call is non-blocking; the scheduler runs in
// its own goroutine managed by robfig/cron.
func (s *Scheduler) Start(cronExpr string, onFire func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.crontab != nil {
		s.crontab.Stop()
	}

	normalized := strings.ReplaceAll(cronExpr, "?", "*")
	c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc(normalized, onFire) // expression is pre-validated by parsing package
	c.Start()
	s.crontab = c
}

// Stop halts the scheduler. Safe to call multiple times.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.crontab != nil {
		s.crontab.Stop()
		s.crontab = nil
	}
}
