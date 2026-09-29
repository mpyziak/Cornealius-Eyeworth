package config

import (
	"encoding/json"
	"os"
)

// Config is the application's persisted settings.
type Config struct {
	CronExpression        string  `json:"CronExpression"`
	StandUpCronExpression string  `json:"StandUpCronExpression"`
	Language              *string `json:"Language"`
}

// Store loads and saves a Config. Repository is its only implementation.
type Store interface {
	Load() (*Config, error)
	Save(*Config) error
}

// Repository is a Store backed by a JSON file on disk.
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
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			if saveErr := r.Save(cfg); saveErr != nil {
				return nil, saveErr
			}
			return cfg, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
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
