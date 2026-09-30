//go:build diagnostics

package devdiag

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	maxFileSize     = 10 * 1024 * 1024 // 10 MB per chunk
	memSamplePeriod = 30 * time.Second
	timestampLayout = "2006-01-02T15-04-05Z"
)

var (
	mu      sync.Mutex
	logDir  string
	current *os.File
	written int64

	closeOnce sync.Once
	stopMem   = make(chan struct{})
)

// Init starts the rotating-file logger and memory sampler, writing log
// files under dir. Call Close on shutdown, else the sampler keeps running.
func Init(dir string) error {
	mu.Lock()
	defer mu.Unlock()

	logDir = dir
	if err := openNewFile(); err != nil {
		return err
	}

	go memoryPoller()
	startSysmon()
	return nil
}

// Close stops the sampler and closes the current log file. Safe to call
// more than once.
func Close() {
	closeOnce.Do(func() {
		stopSysmon()
		close(stopMem)

		mu.Lock()
		defer mu.Unlock()

		if current != nil {
			_ = current.Sync()
			_ = current.Close()
			current = nil
		}
	})
}

// Info logs an informational line, followed by a memory sample.
func Info(format string, args ...any) {
	write("INFO ", format, args...)
	logMemory()
}

// Event logs a notable but non-error occurrence, followed by a memory
// sample.
func Event(format string, args ...any) {
	write("EVENT", format, args...)
	logMemory()
}

// Warn logs a recoverable problem, followed by a memory sample.
func Warn(format string, args ...any) {
	write("WARN ", format, args...)
	logMemory()
}

// Err logs an error, followed by a memory sample. Does not exit.
func Err(format string, args ...any) {
	write("ERROR", format, args...)
	logMemory()
}

func write(level, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s] %s  %s\n",
		time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		level,
		msg,
	)

	mu.Lock()
	defer mu.Unlock()

	if current == nil {
		return
	}

	if written+int64(len(line)) > maxFileSize {
		old := current
		if err := openNewFile(); err == nil {
			_ = old.Sync()
			_ = old.Close()
		} else {
			// Rotation failed; record it in the still-open old file and keep
			// appending there rather than losing every line until the next
			// successful rotation.
			_, _ = fmt.Fprintf(old, "[%s] ERROR rotate log file: %v\n", // nothing further to report to if this write itself fails
				time.Now().UTC().Format("2006-01-02T15:04:05.000Z"), err)
		}
	}

	n, _ := fmt.Fprint(current, line)
	written += int64(n)
}

// Call with mu held.
func openNewFile() error {
	name := fmt.Sprintf("cornealius-%s.log", time.Now().UTC().Format(timestampLayout))
	path := filepath.Join(logDir, name)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	current = f
	written = 0
	return nil
}

func memoryPoller() {
	ticker := time.NewTicker(memSamplePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-stopMem:
			return
		case <-ticker.C:
			logMemory()
		}
	}
}

// Previous sample, for deltas.
var prevMemStats struct {
	numGC        uint32
	pauseTotalNs uint64
}

func logMemory() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	goroutines := runtime.NumGoroutine()

	mu.Lock()
	gcDelta := ms.NumGC - prevMemStats.numGC
	pauseDelta := ms.PauseTotalNs - prevMemStats.pauseTotalNs
	prevMemStats.numGC = ms.NumGC
	prevMemStats.pauseTotalNs = ms.PauseTotalNs
	mu.Unlock()

	write("MEM  ",
		"HeapAlloc=%s HeapSys=%s HeapIdle=%s HeapInuse=%s HeapReleased=%s Sys=%s "+
			"NumGC=%d(+%d) GCPause=%s(+%s) Goroutines=%d",
		fmtBytes(ms.HeapAlloc),
		fmtBytes(ms.HeapSys),
		fmtBytes(ms.HeapIdle),
		fmtBytes(ms.HeapInuse),
		fmtBytes(ms.HeapReleased),
		fmtBytes(ms.Sys),
		ms.NumGC, gcDelta,
		fmtDuration(ms.PauseTotalNs), fmtDuration(pauseDelta),
		goroutines,
	)
}

func fmtDuration(ns uint64) string {
	switch {
	case ns >= 1e9:
		return fmt.Sprintf("%.2fs", float64(ns)/1e9)
	case ns >= 1e6:
		return fmt.Sprintf("%.2fms", float64(ns)/1e6)
	case ns >= 1e3:
		return fmt.Sprintf("%.2fµs", float64(ns)/1e3)
	default:
		return fmt.Sprintf("%dns", ns)
	}
}

func fmtBytes(b uint64) string {
	switch {
	case b >= 1024*1024*1024:
		return fmt.Sprintf("%.2fGB", float64(b)/float64(1024*1024*1024))
	case b >= 1024*1024:
		return fmt.Sprintf("%.2fMB", float64(b)/float64(1024*1024))
	case b >= 1024:
		return fmt.Sprintf("%.2fKB", float64(b)/float64(1024))
	default:
		return fmt.Sprintf("%dB", b)
	}
}
