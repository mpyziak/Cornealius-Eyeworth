// Package diagnostics is the always-present logging interface. impl.go
// (behind the diagnostics build tag) forwards to dev-diagnostics; noop.go
// (the default) makes every call a no-op in release builds.
package diagnostics
