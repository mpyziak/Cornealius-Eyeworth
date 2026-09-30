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
	got, err := repo.Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	want := DefaultConfig()
	if got.CronExpression != want.CronExpression {
		t.Errorf("Load().CronExpression = %q, want %q", got.CronExpression, want.CronExpression)
	}
	if got.StandUpCronExpression != want.StandUpCronExpression {
		t.Errorf("Load().StandUpCronExpression = %q, want %q", got.StandUpCronExpression, want.StandUpCronExpression)
	}
	if got.Language != want.Language {
		t.Errorf("Load().Language = %v, want %v", got.Language, want.Language)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to exist after Load(), got error: %v", err)
	}
}
