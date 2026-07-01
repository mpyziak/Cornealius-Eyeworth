//go:build diagnostics

// Package diagnostics is the optional dev-diagnostics implementation.
// It provides structured, rotating file logging and a periodic memory sampler.
// Compiled only when the diagnostics build tag is present; the production
// binary uses the zero-cost stubs in diagnostics/noop.go instead.
//
// Log files are written next to the executable, named by UTC timestamp:
//
//	cornealius-2026-06-22T14-30-00Z.log
//
// When the current file reaches 10 MB a new file is opened automatically.
// Memory usage is sampled every 30 seconds.
package diagnostics

import (
	"fmt"
	"io"
	"log"
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
	logger  *log.Logger

	stopMem = make(chan struct{})
)

// Init opens the first log file in dir and starts the memory sampler.
// dir is typically the directory of the executable.
// Call Close() on shutdown to flush and stop the sampler.
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

// Close flushes and closes the current log file and stops background goroutines.
func Close() {
	stopSysmon()
	close(stopMem)

	mu.Lock()
	defer mu.Unlock()

	if current != nil {
		_ = current.Sync()
		_ = current.Close()
		current = nil
	}
}

// Info logs a plain informational message.
func Info(format string, args ...any) {
	write("INFO ", format, args...)
	logMemory()
}

// Event logs a significant application event (cron fire, notification sent, etc.).
func Event(format string, args ...any) {
	write("EVENT", format, args...)
	logMemory()
}

// Warn logs a warning.
func Warn(format string, args ...any) {
	write("WARN ", format, args...)
	logMemory()
}

// Err logs an error (does not terminate the program).
func Err(format string, args ...any) {
	write("ERROR", format, args...)
	logMemory()
}

// ---- internal ---------------------------------------------------------------

func write(level, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("[%s] %s  %s\n",
		time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		level,
		msg,
	)

	mu.Lock()
	defer mu.Unlock()

	if logger == nil {
		return
	}

	if written+int64(len(line)) > maxFileSize {
		_ = current.Sync()
		_ = current.Close()
		if err := openNewFile(); err != nil {
			// Can't open new chunk — fall through and keep writing to old file
			// (openNewFile already logged to stderr).
		}
	}

	n, _ := fmt.Fprint(current, line)
	written += int64(n)
	_ = logger // suppress unused warning; logger is kept for its prefix/flag behaviour
}

// openNewFile creates a new timestamped log file. Must be called with mu held.
func openNewFile() error {
	name := fmt.Sprintf("cornealius-%s.log", time.Now().UTC().Format(timestampLayout))
	path := filepath.Join(logDir, name)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dev-diagnostics: cannot open log file %s: %v\n", path, err)
		return err
	}

	current = f
	written = 0
	logger = log.New(io.Discard, "", 0) // placeholder; actual writing done via fmt.Fprint
	return nil
}

// memoryPoller samples runtime memory stats and writes them to the log every 30 s.
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

// prevMemStats holds values from the previous sample for delta calculations.
var prevMemStats struct {
	numGC        uint32
	pauseTotalNs uint64
}

func logMemory() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	goroutines := runtime.NumGoroutine()
	gcDelta := ms.NumGC - prevMemStats.numGC
	pauseDelta := ms.PauseTotalNs - prevMemStats.pauseTotalNs
	prevMemStats.numGC = ms.NumGC
	prevMemStats.pauseTotalNs = ms.PauseTotalNs

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
