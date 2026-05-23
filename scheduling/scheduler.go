package scheduling

import (
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

// Runner is the schedule execution contract.
type Runner interface {
	Start(cronExpr string, onFire func())
	Stop()
}

const (
	ReminderTypeEye     = "eye"
	ReminderTypeStandup = "standup"
)

// ScheduleSpec describes the active schedule expressions that should be
// running. A blank StandUpCron expression means the stand-up schedule is off.
type ScheduleSpec struct {
	EyeCron     string
	StandUpCron string
}

// ApplySchedule starts or stops the provided runners based on spec.
// This centralizes scheduler startup logic and keeps UI code focused on
// behaviour rather than cron wiring.
func ApplySchedule(eye Runner, standup Runner, spec ScheduleSpec, onFire func(notificationCategory string)) {
	if spec.EyeCron != "" {
		eye.Start(spec.EyeCron, func() { onFire(ReminderTypeEye) })
	} else {
		eye.Stop()
	}

	if spec.StandUpCron != "" {
		standup.Start(spec.StandUpCron, func() { onFire(ReminderTypeStandup) })
	} else {
		standup.Stop()
	}
}

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
