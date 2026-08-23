// Package main is the entrypoint for the lumibot Discord bot application.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AobaIwaki123/lumibot/pkg/config"
	"github.com/AobaIwaki123/lumibot/pkg/store"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	slog.Info("starting lumibot", "version", version, "commit", commit, "build_date", date)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	if cfg.LogLevel == "debug" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	st, err := store.New(cfg.DatabasePath)
	if err != nil {
		slog.Error("failed to initialize store", "error", err, "path", cfg.DatabasePath)
		os.Exit(1)
	}
	defer func() { _ = st.Close() }()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("lumibot initialization complete", "db_path", cfg.DatabasePath, "api_url", cfg.LumitreeAPIURL)

	<-ctx.Done()
	slog.Info("shutting down lumibot gracefully")
	fmt.Println("lumibot stopped.")
}
