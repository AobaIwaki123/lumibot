// Package main is the entry point for the lumibot application.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AobaIwaki123/lumibot/pkg/bot"
	"github.com/AobaIwaki123/lumibot/pkg/client"
	"github.com/AobaIwaki123/lumibot/pkg/config"
	"github.com/AobaIwaki123/lumibot/pkg/cron"
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
	defer func() {
		_ = st.Close()
	}()

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

	// Initialize and Start Cron Scheduler
	cronJob, err := cron.New(discordBot.Session, st, apiClient)
	if err != nil {
		slog.Error("Failed to initialize cron", "error", err)
		os.Exit(1)
	}
	cronJob.Start()

	// Start Health Check Server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		// Example readiness logic: check DB connection or Discord session
		if discordBot.Session == nil {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("Starting health check server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Health check server failed", "error", err)
		}
	}()

	slog.Info("Lumibot is now running. Press CTRL-C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	slog.Info("Shutting down...")

	// Graceful shutdown of HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Failed to gracefully shutdown health server", "error", err)
	}

	cronJob.Stop()
	_ = discordBot.Stop()
}
