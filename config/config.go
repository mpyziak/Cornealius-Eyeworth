package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CronExpression        string  `json:"CronExpression"`
	StandUpCronExpression string  `json:"StandUpCronExpression"`
	Language              *string `json:"Language"`
}

type Store interface {
	Load() (*Config, error)
	Save(*Config) error
}

type Repository struct {
	path string
}

func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

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

func (r *Repository) Save(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0600)
}
