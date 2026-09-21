package scheduling

import (
	"fmt"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

type Runner interface {
	Start(cronExpr string, onFire func()) error
	Stop()
}

const (
	ReminderTypeEye     = "eye"
	ReminderTypeStandup = "standup"
)

type ScheduleSpec struct {
	EyeCron     string
	StandUpCron string // blank turns the stand-up schedule off
}

// Returns the first rejected expression, and leaves that runner stopped.
// Callers must surface or recover from it.
func ApplySchedule(eye Runner, standup Runner, spec ScheduleSpec, onFire func(notificationCategory string)) error {
	var firstErr error

	if spec.EyeCron != "" {
		if err := eye.Start(spec.EyeCron, func() { onFire(ReminderTypeEye) }); err != nil {
			firstErr = fmt.Errorf("eye schedule %q: %w", spec.EyeCron, err)
			eye.Stop()
		}
	} else {
		eye.Stop()
	}

	if spec.StandUpCron != "" {
		if err := standup.Start(spec.StandUpCron, func() { onFire(ReminderTypeStandup) }); err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("stand-up schedule %q: %w", spec.StandUpCron, err)
			}
			standup.Stop()
		}
	} else {
		standup.Stop()
	}

	return firstErr
}

type Scheduler struct {
	mu      sync.Mutex
	crontab *cron.Cron
}

// Replaces any running schedule. Validates here rather than trusting the
// parsing package - config.json expressions never go through it.
func (s *Scheduler) Start(cronExpr string, onFire func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.crontab != nil {
		s.crontab.Stop()
		s.crontab = nil
	}

	normalized := strings.ReplaceAll(cronExpr, "?", "*")
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc(normalized, onFire); err != nil {
		return err
	}
	c.Start()
	s.crontab = c
	return nil
}

// Safe to call repeatedly.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.crontab != nil {
		s.crontab.Stop()
		s.crontab = nil
	}
}
