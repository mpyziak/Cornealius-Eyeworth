package scheduling

import (
	"fmt"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

// Runner starts and stops a single cron schedule. Scheduler is its only
// implementation.
type Runner interface {
	Start(cronExpr string, onFire func()) error
	Stop()
}

// Category identifies which schedule fired: CategoryEye or CategoryStandUp.
type Category string

const (
	CategoryEye     Category = "eye"
	CategoryStandUp Category = "standup"
)

// Spec is the pair of cron expressions ApplySchedule installs.
type Spec struct {
	EyeCron     string
	StandUpCron string // blank turns the stand-up schedule off
}

// ApplySchedule starts or stops eye and standup according to spec. It
// returns the first rejected expression and leaves that runner stopped;
// callers must surface or recover from it.
func ApplySchedule(eye Runner, standup Runner, spec Spec, onFire func(category Category)) error {
	var firstErr error

	if spec.EyeCron != "" {
		if err := eye.Start(spec.EyeCron, func() { onFire(CategoryEye) }); err != nil {
			firstErr = fmt.Errorf("eye schedule %q: %w", spec.EyeCron, err)
			eye.Stop()
		}
	} else {
		eye.Stop()
	}

	if spec.StandUpCron != "" {
		if err := standup.Start(spec.StandUpCron, func() { onFire(CategoryStandUp) }); err != nil {
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

// Scheduler is a Runner backed by a robfig/cron schedule.
type Scheduler struct {
	mu      sync.Mutex
	crontab *cron.Cron
}

// Start replaces any running schedule with cronExpr, calling onFire each
// time it triggers. Validates here rather than trusting the parsing
// package - config.json expressions never go through it.
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

// Stop stops the running schedule, if any. Safe to call repeatedly.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.crontab != nil {
		s.crontab.Stop()
		s.crontab = nil
	}
}
