package config

import (
	"time"
	"os"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AgentID   string        `yaml:"agent_id"`
	ServerURL string        `yaml:"server_url"`
	Interval  time.Duration `yaml:"interval"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open cofnig: %w", err)
	}
	defer f.Close()
	
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	
	// Env overrides
	if v := os.Getenv("PULSE_AGENT_ID"); v != "" {
		cfg.AgentID = v
	}
	if v := os.Getenv("PULSE_SERVER_URL"); v != "" {
		cfg.ServerURL = v
	}
	if v := os.Getenv("PULSE_INTERVAL"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("parse PULSE_INTERVAL: %w", err)
		}
		cfg.Interval = d
	}

	return &cfg, nil
}