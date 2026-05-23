package scheduling

import (
	"sync"
	"time"
)

// Reminder describes a single reminder event with a type identifier and message.
type Reminder struct {
	NotificationCategory string // "eye" or "standup"
	Message              string
}

// ReminderAggregator is the contract for buffering and flushing reminders.
// Implements Single Responsibility: aggregates similar events within a time window.
type ReminderAggregator interface {
	Add(reminder Reminder)
	Close()
}

// Buffer accumulates reminders fired within a time window and flushes them together.
// Depends on an abstract Flusher, not a concrete notification implementation (Dependency Inversion).
type Buffer struct {
	mu        sync.Mutex
	pending   []Reminder
	windowLen time.Duration
	timer     *time.Timer
	flusher   Flusher
}

// Flusher is the interface for notification delivery (Open/Closed: can add new implementations).
type Flusher interface {
	Flush(reminders []Reminder)
}

// NewBuffer creates a Buffer with the given time window and flusher.
func NewBuffer(windowLen time.Duration, flusher Flusher) *Buffer {
	return &Buffer{
		windowLen: windowLen,
		flusher:   flusher,
	}
}

// Add appends a reminder and starts the flush window if needed.
func (b *Buffer) Add(reminder Reminder) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.pending = append(b.pending, reminder)

	if b.timer == nil {
		// First event in window: start timer
		b.timer = time.AfterFunc(b.windowLen, b.flush)
	}
}

// Close ensures any pending reminders are flushed immediately.
func (b *Buffer) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}
	b.flushLocked()
}

// flush is called by the timer when the window expires.
func (b *Buffer) flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flushLocked()
}

// flushLocked sends pending reminders to the flusher. Must be called with lock held.
func (b *Buffer) flushLocked() {
	if len(b.pending) == 0 {
		return
	}
	pending := b.pending
	b.pending = []Reminder{}
	b.timer = nil

	// Release lock before calling flusher to avoid deadlocks
	b.mu.Unlock()
	b.flusher.Flush(pending)
	b.mu.Lock()
}
