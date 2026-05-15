package config

import (
	"encoding/json"
	"os"
)

// Config holds the user-configurable application settings.
type Config struct {
	CronExpression        string  `json:"CronExpression"`
	StandUpCronExpression string  `json:"StandUpCronExpression"`
	Language              *string `json:"Language"`
}

// Store is the read/write contract for application settings.
type Store interface {
	Load() (*Config, error)
	Save(*Config) error
}

// Repository reads and writes Config as JSON next to the executable.
type Repository struct {
	path string
}

// NewRepository creates a Repository that persists config at path.
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

// Load reads the config file and returns a parsed Config.
func (r *Repository) Load() (*Config, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes cfg to the config file as indented JSON.
func (r *Repository) Save(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0600)
}
