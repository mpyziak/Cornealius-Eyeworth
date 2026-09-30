// Package devdiag is the real logger and memory sampler behind the
// diagnostics build tag: rotating log files next to the exe, plus (with the
// sysmon tag too) OS-level event monitoring on Windows. Release builds
// compile github.com/mpyziak/cornealius-eyeworth/diagnostics's noop.go
// instead.
package devdiag
