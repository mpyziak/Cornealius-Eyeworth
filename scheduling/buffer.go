package scheduling

import (
	"sync"
	"time"
)

type Reminder struct {
	NotificationCategory string // "eye" or "standup"
	Message              string
}

type ReminderAggregator interface {
	Add(reminder Reminder)
	Close()
}

type Buffer struct {
	mu        sync.Mutex
	pending   []Reminder
	windowLen time.Duration
	timer     *time.Timer
	flusher   Flusher
}

type Flusher interface {
	Flush(reminders []Reminder)
}

func NewBuffer(windowLen time.Duration, flusher Flusher) *Buffer {
	return &Buffer{
		windowLen: windowLen,
		flusher:   flusher,
	}
}

func (b *Buffer) Add(reminder Reminder) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.pending = append(b.pending, reminder)

	if b.timer == nil {
		b.timer = time.AfterFunc(b.windowLen, b.flush)
	}
}

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
