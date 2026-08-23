// Package bot implements the Discord bot logic and slash command handlers.
package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/AobaIwaki123/lumibot/pkg/client"
	"github.com/AobaIwaki123/lumibot/pkg/store"
)

// Bot represents the Discord bot instance.
type Bot struct {
	Session *discordgo.Session
	Store   store.Store
	Client  client.Client
	AppID   string
}

// New creates a new Bot instance.
func New(token string, s store.Store, c client.Client) (*Bot, error) {
	if !strings.HasPrefix(token, "Bot ") {
		token = "Bot " + token
	}

	session, err := discordgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create discord session: %w", err)
	}

	bot := &Bot{
		Session: session,
		Store:   s,
		Client:  c,
	}

	// Register handlers
	bot.Session.AddHandler(bot.onReady)
	bot.Session.AddHandler(bot.onInteractionCreate)

	return bot, nil
}

// Start opens the websocket connection to Discord and registers slash commands.
func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	slog.Info("Discord session opened successfully")

	err := b.registerCommands()
	if err != nil {
		return fmt.Errorf("failed to register commands: %w", err)
	}

	return nil
}

// Stop gracefully closes the websocket connection.
func (b *Bot) Stop() error {
	return b.Session.Close()
}

func (b *Bot) onReady(s *discordgo.Session, r *discordgo.Ready) {
	slog.Info("Bot is ready", "username", s.State.User.Username, "discriminator", s.State.User.Discriminator)
	b.AppID = s.State.User.ID
}
