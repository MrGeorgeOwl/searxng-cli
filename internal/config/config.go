package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

const appName = "searxng-cli"

// Config is the on-disk configuration for searxng-cli.
type Config struct {
	Environments map[string]Environment `json:"environments"`
}

// Environment describes a SearXNG instance.
type Environment struct {
	URL string `json:"url"`
}

// DefaultPath returns the default config path.
func DefaultPath() string {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, appName, "config.json")
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil || userConfigDir == "" {
		return filepath.Join(".", "config.json")
	}

	return filepath.Join(userConfigDir, appName, "config.json")
}

// Load reads and validates configuration from path.
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("load config %q: %w", path, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks whether the configuration is usable.
func (c Config) Validate() error {
	if len(c.Environments) == 0 {
		return errors.New("config must define environments")
	}
	if _, ok := c.Environments["default"]; !ok {
		return errors.New("config must define default environment")
	}

	for name, env := range c.Environments {
		if env.URL == "" {
			return fmt.Errorf("environment %q must define url", name)
		}
		parsed, err := url.Parse(env.URL)
		if err != nil {
			return fmt.Errorf("environment %q has invalid url: %w", name, err)
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return fmt.Errorf("environment %q url must use http or https", name)
		}
		if parsed.Host == "" {
			return fmt.Errorf("environment %q url must include a host", name)
		}
	}

	return nil
}

// Environment returns a named environment from the config.
func (c Config) Environment(name string) (Environment, error) {
	env, ok := c.Environments[name]
	if !ok {
		return Environment{}, fmt.Errorf("environment %q is not defined", name)
	}
	return env, nil
}
