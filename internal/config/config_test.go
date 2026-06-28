package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
  "environments": {
    "default": {"url": "http://localhost:8080"},
    "other": {"url": "https://example.com"}
  }
}`

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	env, err := cfg.Environment("other")
	if err != nil {
		t.Fatalf("Environment() error = %v", err)
	}
	if env.URL != "https://example.com" {
		t.Fatalf("env.URL = %q", env.URL)
	}
}

func TestLinuxDefaultPathUsesXDGConfigHome(t *testing.T) {
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)

	got, err := linuxDefaultPath()
	if err != nil {
		t.Fatalf("linuxDefaultPath() error = %v", err)
	}

	want := filepath.Join(configHome, appName, "config.json")
	if got != want {
		t.Fatalf("linuxDefaultPath() = %q, want %q", got, want)
	}
}

func TestLinuxDefaultPathFallsBackToHomeDotConfig(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", homeDir)

	got, err := linuxDefaultPath()
	if err != nil {
		t.Fatalf("linuxDefaultPath() error = %v", err)
	}

	want := filepath.Join(homeDir, ".config", appName, "config.json")
	if got != want {
		t.Fatalf("linuxDefaultPath() = %q, want %q", got, want)
	}
}

func TestMacOSDefaultPathUsesHomeDotConfig(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("HOME", homeDir)

	got, err := macOSDefaultPath()
	if err != nil {
		t.Fatalf("macOSDefaultPath() error = %v", err)
	}

	want := filepath.Join(homeDir, ".config", appName, "config.json")
	if got != want {
		t.Fatalf("macOSDefaultPath() = %q, want %q", got, want)
	}
}

func TestValidateRequiresDefaultEnvironment(t *testing.T) {
	cfg := Config{Environments: map[string]Environment{
		"other": {URL: "http://localhost:8080"},
	}}

	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}

func TestValidateRejectsInvalidURL(t *testing.T) {
	cfg := Config{Environments: map[string]Environment{
		"default": {URL: "localhost:8080"},
	}}

	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}
