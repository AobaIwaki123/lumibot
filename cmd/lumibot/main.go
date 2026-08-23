package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AobaIwaki123/lumibot/pkg/bot"
	"github.com/AobaIwaki123/lumibot/pkg/client"
	"github.com/AobaIwaki123/lumibot/pkg/config"
	"github.com/AobaIwaki123/lumibot/pkg/store"
)

func main() {
	cfg := config.Load()

	// Initialize Store
	st, err := store.New(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer st.Close()

	// Initialize API Client
	apiClient, err := client.NewLumitreeClient(cfg.LumitreeAPIURL)
	if err != nil {
		slog.Error("Failed to initialize lumitree client", "error", err)
		os.Exit(1)
	}

	// Initialize and Start Discord Bot
	discordBot, err := bot.New(cfg.DiscordToken, st, apiClient)
	if err != nil {
		slog.Error("Failed to initialize bot", "error", err)
		os.Exit(1)
	}

	if err := discordBot.Start(); err != nil {
		slog.Error("Failed to start bot", "error", err)
		os.Exit(1)
	}

	slog.Info("Lumibot is now running. Press CTRL-C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	slog.Info("Shutting down...")
	_ = discordBot.Stop()
}
