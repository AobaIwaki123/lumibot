package config_test

import (
	"os"
	"testing"

	"github.com/AobaIwaki123/lumibot/pkg/config"
)

func TestLoad_Defaults(t *testing.T) {
	_ = os.Unsetenv("DISCORD_TOKEN")
	_ = os.Unsetenv("DISCORD_APP_ID")
	_ = os.Unsetenv("LUMITREE_API_URL")
	_ = os.Unsetenv("SQLITE_DB_PATH")
	_ = os.Unsetenv("LOG_LEVEL")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.LumitreeAPIURL != "https://lumitree.aooba.net" {
		t.Errorf("expected default LumitreeAPIURL 'https://lumitree.aooba.net', got '%s'", cfg.LumitreeAPIURL)
	}
	if cfg.DatabasePath != "lumibot.db" {
		t.Errorf("expected default DatabasePath 'lumibot.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default LogLevel 'info', got '%s'", cfg.LogLevel)
	}

	if err := cfg.Validate(); err == nil {
		t.Errorf("expected error when DISCORD_TOKEN is missing, got nil")
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("DISCORD_TOKEN", "mock-token")
	t.Setenv("DISCORD_APP_ID", "123456789")
	t.Setenv("LUMITREE_API_URL", "http://localhost:8080/")
	t.Setenv("SQLITE_DB_PATH", "/tmp/test.db")
	t.Setenv("LOG_LEVEL", "DEBUG")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.DiscordToken != "mock-token" {
		t.Errorf("expected DiscordToken 'mock-token', got '%s'", cfg.DiscordToken)
	}
	if cfg.DiscordAppID != "123456789" {
		t.Errorf("expected DiscordAppID '123456789', got '%s'", cfg.DiscordAppID)
	}
	if cfg.LumitreeAPIURL != "http://localhost:8080" {
		t.Errorf("expected trimmed LumitreeAPIURL 'http://localhost:8080', got '%s'", cfg.LumitreeAPIURL)
	}
	if cfg.DatabasePath != "/tmp/test.db" {
		t.Errorf("expected DatabasePath '/tmp/test.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel 'debug', got '%s'", cfg.LogLevel)
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got error: %v", err)
	}
}
