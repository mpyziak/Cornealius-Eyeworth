//go:build diagnostics

// Rotating log files next to the exe
// Release builds get diagnostics/noop.go instead.
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

// Close() on shutdown, else the sampler keeps running.
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

func Info(format string, args ...any) {
	write("INFO ", format, args...)
	logMemory()
}

func Event(format string, args ...any) {
	write("EVENT", format, args...)
	logMemory()
}

func Warn(format string, args ...any) {
	write("WARN ", format, args...)
	logMemory()
}

// Does not exit.
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

	if logger == nil {
		return
	}

	if written+int64(len(line)) > maxFileSize {
		_ = current.Sync()
		_ = current.Close()
		if err := openNewFile(); err != nil {
			// Keep writing to the old file; openNewFile already told stderr.
		}
	}

	n, _ := fmt.Fprint(current, line)
	written += int64(n)
	_ = logger // suppress unused warning; logger is kept for its prefix/flag behaviour
}

// Call with mu held.
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
