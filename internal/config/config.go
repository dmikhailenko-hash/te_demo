package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token,omitempty"`
}

const DefaultBaseURL = "https://webhooks-clientapi.traderevolution.com/traderevolution/v1"

func cfgPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".te_demo")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func Load() (*Config, error) {
	p, err := cfgPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &Config{BaseURL: DefaultBaseURL}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	p, err := cfgPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token not set — run: te_demo config --token YOUR_TOKEN")
	}
	return nil
}
