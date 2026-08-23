// Package config provides application configuration loaded from environment variables.
package config

import (
	"fmt"
	"os"
	"strings"
)

// Config represents the application configuration.
type Config struct {
	DiscordToken   string
	DiscordAppID   string
	LumitreeAPIURL string
	DatabasePath   string
	LogLevel       string
}

// Load loads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	token := strings.TrimSpace(os.Getenv("DISCORD_TOKEN"))
	appID := strings.TrimSpace(os.Getenv("DISCORD_APP_ID"))

	apiURL := strings.TrimSpace(os.Getenv("LUMITREE_API_URL"))
	if apiURL == "" {
		apiURL = "https://lumitree.aooba.net"
	}
	apiURL = strings.TrimRight(apiURL, "/")

	dbPath := strings.TrimSpace(os.Getenv("SQLITE_DB_PATH"))
	if dbPath == "" {
		dbPath = "lumibot.db"
	}

	logLevel := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		DiscordToken:   token,
		DiscordAppID:   appID,
		LumitreeAPIURL: apiURL,
		DatabasePath:   dbPath,
		LogLevel:       logLevel,
	}, nil
}

// Validate checks whether required configuration values are present.
func (c *Config) Validate() error {
	if c.DiscordToken == "" {
		return fmt.Errorf("DISCORD_TOKEN is required")
	}
	return nil
}
