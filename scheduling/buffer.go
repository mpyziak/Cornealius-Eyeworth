package scheduling

import (
	"sync"
	"time"
)

// Reminder is one eye-care or stand-up notification waiting to be flushed.
type Reminder struct {
	NotificationCategory string // "eye" or "standup"
	Message              string
}

// ReminderAggregator collects reminders and flushes them as a batch.
// Buffer is its only implementation.
type ReminderAggregator interface {
	Add(reminder Reminder)
	Close()
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
	defer b.mu.Unlock()

	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}
	b.flushLocked()
}

func (b *Buffer) flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.flushLocked()
}

func (b *Buffer) flushLocked() {
	if len(b.pending) == 0 {
		return
	}
	pending := b.pending
	b.pending = []Reminder{}
	b.timer = nil

	b.mu.Unlock()
	b.flusher.Flush(pending)
	b.mu.Lock()
}
