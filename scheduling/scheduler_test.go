package scheduling

import (
	"errors"
	"testing"
)

type fakeRunner struct {
	started   int
	stopped   int
	lastExpr  string
	lastFire  func()
	startErr  error
	startErrN int // return startErr on the nth Start (1-based); 0 = always
}

func (f *fakeRunner) Start(cronExpr string, onFire func()) error {
	f.started++
	f.lastExpr = cronExpr
	f.lastFire = onFire
	if f.startErr != nil && (f.startErrN == 0 || f.startErrN == f.started) {
		return f.startErr
	}
	return nil
}

func (f *fakeRunner) Stop() { f.stopped++ }

func TestApplyScheduleStartsBoth(t *testing.T) {
	eye, standup := &fakeRunner{}, &fakeRunner{}

	err := ApplySchedule(eye, standup,
		Spec{EyeCron: "0 20 * * * ?", StandUpCron: "0 55 * * * ?"},
		func(Category) {})
	if err != nil {
		t.Fatalf("ApplySchedule returned %v, want nil", err)
	}

	if eye.started != 1 || eye.lastExpr != "0 20 * * * ?" {
		t.Errorf("eye: started=%d expr=%q, want started=1 expr=%q", eye.started, eye.lastExpr, "0 20 * * * ?")
	}
	if standup.started != 1 || standup.lastExpr != "0 55 * * * ?" {
		t.Errorf("standup: started=%d expr=%q, want started=1 expr=%q", standup.started, standup.lastExpr, "0 55 * * * ?")
	}
}

func TestApplyScheduleStopsOnEmptyExpression(t *testing.T) {
	eye, standup := &fakeRunner{}, &fakeRunner{}

	if err := ApplySchedule(eye, standup,
		Spec{EyeCron: "0 20 * * * ?", StandUpCron: ""},
		func(Category) {}); err != nil {
		t.Fatalf("ApplySchedule returned %v, want nil", err)
	}

	if standup.started != 0 {
		t.Errorf("standup started %d times for an empty expression, want 0", standup.started)
	}
	if standup.stopped != 1 {
		t.Errorf("standup stopped %d times, want 1", standup.stopped)
	}
}

func TestApplyScheduleStopsOnEmptyEyeExpression(t *testing.T) {
	eye, standup := &fakeRunner{}, &fakeRunner{}

	if err := ApplySchedule(eye, standup,
		Spec{EyeCron: "", StandUpCron: "0 55 * * * ?"},
		func(Category) {}); err != nil {
		t.Fatalf("ApplySchedule returned %v, want nil", err)
	}

	if eye.started != 0 {
		t.Errorf("eye started %d times for an empty expression, want 0", eye.started)
	}
	if eye.stopped != 1 {
		t.Errorf("eye stopped %d times, want 1", eye.stopped)
	}
}

func TestApplyScheduleReportsFirstErrorWhenBothFail(t *testing.T) {
	eyeErr := errors.New("bad eye expression")
	standupErr := errors.New("bad standup expression")
	eye := &fakeRunner{startErr: eyeErr}
	standup := &fakeRunner{startErr: standupErr}

	err := ApplySchedule(eye, standup,
		Spec{EyeCron: "nonsense", StandUpCron: "also nonsense"},
		func(Category) {})

	if err == nil {
		t.Fatal("ApplySchedule returned nil, want the eye runner's error")
	}
	if !errors.Is(err, eyeErr) {
		t.Errorf("error %v does not wrap the eye runner's error (source order), want it first", err)
	}
	if errors.Is(err, standupErr) {
		t.Errorf("error %v wraps the standup runner's error, want only the first (eye)", err)
	}
	if eye.stopped != 1 {
		t.Errorf("failed eye runner stopped %d times, want 1", eye.stopped)
	}
	if standup.stopped != 1 {
		t.Errorf("failed standup runner stopped %d times, want 1", standup.stopped)
	}
}

// Stopped, not running with no entries - that looks like a working app.
func TestApplyScheduleReportsAndStopsOnStartError(t *testing.T) {
	boom := errors.New("bad expression")
	eye := &fakeRunner{startErr: boom}
	standup := &fakeRunner{}

	err := ApplySchedule(eye, standup,
		Spec{EyeCron: "nonsense", StandUpCron: "0 55 * * * ?"},
		func(Category) {})

	if err == nil {
		t.Fatal("ApplySchedule returned nil, want the runner's error")
	}
	if !errors.Is(err, boom) {
		t.Errorf("error %v does not wrap the runner's error", err)
	}
	if eye.stopped != 1 {
		t.Errorf("failed eye runner stopped %d times, want 1", eye.stopped)
	}
	if standup.started != 1 {
		t.Errorf("standup started %d times, want 1", standup.started)
	}
}

func TestApplyScheduleFiresCorrectCategory(t *testing.T) {
	eye, standup := &fakeRunner{}, &fakeRunner{}

	var got []Category
	if err := ApplySchedule(eye, standup,
		Spec{EyeCron: "0 20 * * * ?", StandUpCron: "0 55 * * * ?"},
		func(c Category) { got = append(got, c) }); err != nil {
		t.Fatalf("ApplySchedule: %v", err)
	}

	eye.lastFire()
	standup.lastFire()

	if len(got) != 2 || got[0] != CategoryEye || got[1] != CategoryStandUp {
		t.Errorf("fired categories = %v, want [%s %s]", got, CategoryEye, CategoryStandUp)
	}
}

func TestSchedulerStartRejectsInvalidExpression(t *testing.T) {
	var s Scheduler
	t.Cleanup(s.Stop)

	if err := s.Start("not a cron expression", func() {}); err == nil {
		t.Fatal("Start accepted an invalid expression, want an error")
	}
	if s.crontab != nil {
		t.Error("Start left a crontab behind after rejecting the expression")
	}
}

func TestSchedulerStartAcceptsQuartzWildcard(t *testing.T) {
	var s Scheduler
	t.Cleanup(s.Stop)

	// Quartz '?', normalised to '*' by Start.
	if err := s.Start("0 20,40,55 * * * ?", func() {}); err != nil {
		t.Fatalf("Start rejected a valid Quartz expression: %v", err)
	}
	if s.crontab == nil {
		t.Error("Start reported success but stored no crontab")
	}
}

// A second Start replaces the first rather than stacking another cron.
func TestSchedulerStartReplacesPreviousSchedule(t *testing.T) {
	var s Scheduler
	t.Cleanup(s.Stop)

	if err := s.Start("0 20 * * * ?", func() {}); err != nil {
		t.Fatalf("first Start: %v", err)
	}
	first := s.crontab

	if err := s.Start("0 30 * * * ?", func() {}); err != nil {
		t.Fatalf("second Start: %v", err)
	}
	if s.crontab == first {
		t.Error("second Start reused the first crontab instead of replacing it")
	}
	if len(s.crontab.Entries()) != 1 {
		t.Errorf("crontab has %d entries after restart, want 1", len(s.crontab.Entries()))
	}
}

func TestSchedulerStopIsIdempotent(t *testing.T) {
	var s Scheduler

	if err := s.Start("0 20 * * * ?", func() {}); err != nil {
		t.Fatalf("Start: %v", err)
	}
	s.Stop()
	s.Stop() // must not panic on an already-stopped scheduler

	if s.crontab != nil {
		t.Error("Stop left a crontab behind")
	}
}
