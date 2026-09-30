package scheduling

import (
	"sync"
	"time"
)

// Reminder is one eye-care or stand-up notification waiting to be flushed.
type Reminder struct {
	Category Category
	Message  string
}

// Buffer collects reminders that fire close together and flushes them as
// one batch after windowLen, instead of one notification per fire.
type Buffer struct {
	mu        sync.Mutex
	pending   []Reminder
	windowLen time.Duration
	timer     *time.Timer
	flusher   Flusher
}

// Flusher delivers a batch of reminders, e.g. as one notification.
type Flusher interface {
	Flush(reminders []Reminder)
}

// NewBuffer returns a Buffer that flushes through flusher windowLen after
// the first reminder in each batch arrives.
func NewBuffer(windowLen time.Duration, flusher Flusher) *Buffer {
	return &Buffer{
		windowLen: windowLen,
		flusher:   flusher,
	}
}

// Add queues reminder, arming the flush timer if it is not already running.
func (b *Buffer) Add(reminder Reminder) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.pending = append(b.pending, reminder)

	if b.timer == nil {
		b.timer = time.AfterFunc(b.windowLen, b.flush)
	}
}

// Close cancels the flush timer and flushes whatever is pending.
func (b *Buffer) Close() {
	b.mu.Lock()
	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}
	pending := b.flushLocked()
	b.mu.Unlock()

	if pending != nil {
		b.flusher.Flush(pending)
	}
}

func (b *Buffer) flush() {
	b.mu.Lock()
	pending := b.flushLocked()
	b.mu.Unlock()

	if pending != nil {
		b.flusher.Flush(pending)
	}
}

// flushLocked clears and returns the pending batch. The caller must hold mu
// and is responsible for calling Flush after releasing it, so a panicking
// Flusher can't leave mu locked by one goroutine and unlocked by another.
func (b *Buffer) flushLocked() []Reminder {
	if len(b.pending) == 0 {
		return nil
	}
	pending := b.pending
	b.pending = nil
	b.timer = nil
	return pending
}
