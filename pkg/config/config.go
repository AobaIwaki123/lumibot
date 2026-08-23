package config

import "os"

// Config holds the application configuration.
type Config struct {
	DBPath         string
	DiscordToken   string
	LumitreeAPIURL string
}

// Load reads the configuration from environment variables and provides defaults.
func Load() *Config {
	dbPath := os.Getenv("LUMIBOT_DB_PATH")
	if dbPath == "" {
		dbPath = "lumibot.db"
	}

	discordToken := os.Getenv("DISCORD_TOKEN")

	apiURL := os.Getenv("LUMITREE_API_URL")
	if apiURL == "" {
		apiURL = "https://api.lumitree.example.com"
	}

	return &Config{
		DBPath:         dbPath,
		DiscordToken:   discordToken,
		LumitreeAPIURL: apiURL,
	}
}
