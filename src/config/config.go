package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type AutoGc struct {
	Enabled       *bool `json:"enabled,omitempty"`
	IntervalHours *int  `json:"intervalHours,omitempty"`
}

type Config struct {
	AutoGc *AutoGc `json:"autoGc,omitempty"`
}

type State struct {
	LastRunAt *time.Time `json:"lastRunAt,omitempty"`
}

const defaultIntervalHours = 6

func GetWutHome() string {
	if home := os.Getenv("WUT_HOME"); home != "" {
		return home
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".wut")
}

func GetConfigPath() string {
	return filepath.Join(GetWutHome(), "config.json")
}

func GetStatePath() string {
	return filepath.Join(GetWutHome(), "state.json")
}

func ReadConfig() Config {
	data, err := os.ReadFile(GetConfigPath())
	if err != nil {
		return Config{}
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}
	}
	return cfg
}

func ReadState() State {
	data, err := os.ReadFile(GetStatePath())
	if err != nil {
		return State{}
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return State{}
	}
	return state
}

func WriteState(state State) error {
	statePath := GetStatePath()
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, 0644)
}

func GetAutoGcSettings() (enabled bool, interval time.Duration) {
	cfg := ReadConfig()

	enabled = true
	if cfg.AutoGc != nil && cfg.AutoGc.Enabled != nil {
		enabled = *cfg.AutoGc.Enabled
	}

	hours := defaultIntervalHours
	if cfg.AutoGc != nil && cfg.AutoGc.IntervalHours != nil {
		hours = *cfg.AutoGc.IntervalHours
		if hours < 0 {
			hours = 0
		}
	}

	return enabled, time.Duration(hours) * time.Hour
}
