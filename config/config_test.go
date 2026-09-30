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

func TestLoad_RoundTripsExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	repo := NewRepository(path)

	lang := "de"
	want := &Config{
		CronExpression:        "0 20 * * * ?",
		StandUpCronExpression: "0 55 * * * ?",
		Language:              &lang,
	}
	if err := repo.Save(want); err != nil {
		t.Fatalf("Save() returned unexpected error: %v", err)
	}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if got.CronExpression != want.CronExpression {
		t.Errorf("Load().CronExpression = %q, want %q", got.CronExpression, want.CronExpression)
	}
	if got.StandUpCronExpression != want.StandUpCronExpression {
		t.Errorf("Load().StandUpCronExpression = %q, want %q", got.StandUpCronExpression, want.StandUpCronExpression)
	}
	if got.Language == nil || *got.Language != *want.Language {
		gotLang := "nil"
		if got.Language != nil {
			gotLang = *got.Language
		}
		t.Errorf("Load().Language = %s, want %q", gotLang, *want.Language)
	}
}

func TestLoad_ReturnsErrorOnCorruptJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("not json"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	repo := NewRepository(path)

	if _, err := repo.Load(); err == nil {
		t.Fatal("Load() on corrupt JSON returned nil error, want a parse error")
	}
}
