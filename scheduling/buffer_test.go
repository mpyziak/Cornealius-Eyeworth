package scheduling

import (
	"sync"
	"testing"
	"time"
)

type recordingFlusher struct {
	mu    sync.Mutex
	calls [][]Reminder
}

func (f *recordingFlusher) Flush(reminders []Reminder) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, reminders)
}

func (f *recordingFlusher) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func TestBufferAddThenClose(t *testing.T) {
	f := &recordingFlusher{}
	b := NewBuffer(time.Hour, f) // long window: Close must flush explicitly, not the timer
	b.Add(Reminder{Category: CategoryEye, Message: "look away"})
	b.Close()

	if got := f.callCount(); got != 1 {
		t.Fatalf("Flush called %d times, want 1", got)
	}
}

type panicFlusher struct{}

func (panicFlusher) Flush([]Reminder) { panic("boom") }

// Regression test: flushLocked used to unlock mu and call Flush while still
// "inside" the caller's locked section, so a panicking Flush left mu
// unlocked and the caller's deferred Unlock panicked a second time
// ("unlock of unlocked mutex"), masking the original panic.
func TestBufferFlushPanicPropagatesCleanly(t *testing.T) {
	b := NewBuffer(time.Hour, panicFlusher{})
	b.Add(Reminder{Category: CategoryEye, Message: "x"})

	defer func() {
		r := recover()
		if r != "boom" {
			t.Fatalf("recover() = %v, want the Flusher's own panic (\"boom\")", r)
		}
	}()
	b.Close()
}

// Exercises the mutex-guarded slice/timer under concurrent Add and Close
// calls so -race actually reaches this file (it previously ran at 0%
// coverage under go test -race, per doc/potential-fixes-smells.md §1.8).
func TestBufferConcurrentAddAndClose(t *testing.T) {
	f := &recordingFlusher{}
	b := NewBuffer(time.Millisecond, f)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Add(Reminder{Category: CategoryStandUp, Message: "y"})
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		b.Close()
	}()
	wg.Wait()

	b.Close() // must not panic or deadlock when called again after the race above
}
