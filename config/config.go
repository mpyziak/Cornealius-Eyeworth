package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// Config is the application's persisted settings.
type Config struct {
	CronExpression        string  `json:"CronExpression"`
	StandUpCronExpression string  `json:"StandUpCronExpression"`
	Language              *string `json:"Language"`
}

// Repository loads and saves a Config, backed by a JSON file on disk.
type Repository struct {
	path string
}

// NewRepository returns a Repository that reads and writes path.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

// Load reads the Config from disk, creating it with DefaultConfig's values
// if it does not yet exist.
func (r *Repository) Load() (*Config, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cfg := DefaultConfig()
			if saveErr := r.Save(cfg); saveErr != nil {
				return nil, saveErr
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("read %s: %w", r.path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", r.path, err)
	}
	return &cfg, nil
}

// Save writes cfg to disk as JSON.
func (r *Repository) Save(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0600)
}
