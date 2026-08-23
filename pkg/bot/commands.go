package bot

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "add",
		Description: "Register a lumitree calendar for daily updates in this server.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "calendar_id",
				Description: "The lumitree calendar ID or alias",
				Required:    true,
			},
		},
	},
	{
		Name:        "list",
		Description: "List all registered lumitree calendars in this server.",
	},
	{
		Name:        "remove",
		Description: "Unregister a lumitree calendar.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "calendar_id",
				Description: "The lumitree calendar ID or alias to remove",
				Required:    true,
			},
		},
	},
	{
		Name:        "today",
		Description: "Fetch and display today's events for all registered calendars.",
	},
}

func (b *Bot) registerCommands() error {
	slog.Info("Registering application commands...")
	// Register globally (takes up to 1h to propagate) or per guild (instant).
	// For this MVP, we register globally (guildID="").
	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("cannot create '%v' command: %w", cmd.Name, err)
		}
	}
	slog.Info("Successfully registered application commands")
	return nil
}
