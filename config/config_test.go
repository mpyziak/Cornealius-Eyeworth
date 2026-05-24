package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCreatesDefaultConfigWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	repo := NewRepository(path)
	cfg, err := repo.Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.CronExpression != "0 20,40,55 * * * ?" {
		t.Fatalf("expected default CronExpression, got %q", cfg.CronExpression)
	}
	if cfg.StandUpCronExpression != "0 55 * * * ?" {
		t.Fatalf("expected default StandUpCronExpression, got %q", cfg.StandUpCronExpression)
	}
	if cfg.Language != nil {
		t.Fatalf("expected default Language nil, got %v", cfg.Language)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to exist after Load(), got error: %v", err)
	}
}
