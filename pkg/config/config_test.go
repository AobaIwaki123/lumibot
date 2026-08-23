package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Unset to ensure defaults
	os.Unsetenv("LUMIBOT_DB_PATH")
	os.Unsetenv("DISCORD_TOKEN")
	os.Unsetenv("LUMITREE_API_URL")

	cfg := Load()

	if cfg.DBPath != "lumibot.db" {
		t.Errorf("expected default DBPath 'lumibot.db', got '%s'", cfg.DBPath)
	}
	if cfg.DiscordToken != "" {
		t.Errorf("expected empty default DiscordToken, got '%s'", cfg.DiscordToken)
	}
	if cfg.LumitreeAPIURL != "https://api.lumitree.example.com" {
		t.Errorf("expected default LumitreeAPIURL, got '%s'", cfg.LumitreeAPIURL)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	os.Setenv("LUMIBOT_DB_PATH", "/tmp/test.db")
	os.Setenv("DISCORD_TOKEN", "test-token")
	os.Setenv("LUMITREE_API_URL", "http://localhost:8080")
	defer os.Unsetenv("LUMIBOT_DB_PATH")
	defer os.Unsetenv("DISCORD_TOKEN")
	defer os.Unsetenv("LUMITREE_API_URL")

	cfg := Load()

	if cfg.DBPath != "/tmp/test.db" {
		t.Errorf("expected DBPath '/tmp/test.db', got '%s'", cfg.DBPath)
	}
	if cfg.DiscordToken != "test-token" {
		t.Errorf("expected DiscordToken 'test-token', got '%s'", cfg.DiscordToken)
	}
	if cfg.LumitreeAPIURL != "http://localhost:8080" {
		t.Errorf("expected LumitreeAPIURL 'http://localhost:8080', got '%s'", cfg.LumitreeAPIURL)
	}
}
