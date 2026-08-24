// Package cron provides scheduled background jobs for the bot.
package cron

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AobaIwaki123/lumibot/pkg/client"
	"github.com/AobaIwaki123/lumibot/pkg/store"
	"github.com/bwmarrin/discordgo"
	robfigcron "github.com/robfig/cron/v3"
)

// Cron handles scheduled tasks.
type Cron struct {
	scheduler *robfigcron.Cron
	session   *discordgo.Session
	store     store.Store
	api       client.Client
}

// New creates a new Cron instance.
func New(s *discordgo.Session, st store.Store, api client.Client) (*Cron, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %w", err)
	}

	c := robfigcron.New(robfigcron.WithLocation(loc))

	cronJob := &Cron{
		scheduler: c,
		session:   s,
		store:     st,
		api:       api,
	}

	// MVP: Hardcoded to 08:00 AM JST daily
	_, err = c.AddFunc("0 8 * * *", cronJob.dailyBroadcast)
	if err != nil {
		return nil, fmt.Errorf("failed to add cron func: %w", err)
	}

	return cronJob, nil
}

// Start begins the cron scheduler.
func (c *Cron) Start() {
	c.scheduler.Start()
	slog.Info("Cron scheduler started")
}

// Stop halts the cron scheduler.
func (c *Cron) Stop() {
	c.scheduler.Stop()
	slog.Info("Cron scheduler stopped")
}

func (c *Cron) dailyBroadcast() {
	slog.Info("Running daily broadcast job")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// In MVP, we might need a way to list ALL guilds that have subscriptions,
	// but currently the Store only has ListSubscriptions(guildID).
	// To do this efficiently, we need a new Store method: GetAllGuildIDs()
	// For now, we will add it to the Store interface.

	guildIDs, err := c.store.GetAllGuildIDs(ctx)
	if err != nil {
		slog.Error("Failed to get guild IDs for broadcast", "error", err)
		return
	}

	for _, gID := range guildIDs {
		c.broadcastToGuild(ctx, gID)
	}
}

func (c *Cron) broadcastToGuild(ctx context.Context, guildID string) {
	subs, err := c.store.ListSubscriptions(ctx, guildID)
	if err != nil || len(subs) == 0 {
		return
	}

	gs, err := c.store.GetGuildSettings(ctx, guildID)
	if err != nil || gs.NotifyChannelID == "" {
		slog.Warn("No notify channel set for guild, skipping broadcast", "guild_id", guildID)
		return
	}

	for _, sub := range subs {
		events, err := c.api.GetEvents(ctx, sub.CalendarID, nil) // MVP fetches upcoming/today
		if err != nil {
			slog.Error("Failed to fetch events for cron", "calendar_id", sub.CalendarID, "error", err)
			continue
		}

		if len(events) == 0 {
			continue // Don't spam if no events
		}

		var embeds []*discordgo.MessageEmbed
		for idx, ev := range events {
			if idx >= 5 {
				break
			}
			desc := ""
			if ev.Description != nil {
				desc = *ev.Description
				if len(desc) > 200 {
					desc = desc[:197] + "..."
				}
			}
			loc := "TBD"
			if ev.Location != nil && *ev.Location != "" {
				loc = *ev.Location
			}
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       ev.Title,
				Description: desc,
				Color:       0x00b0f4,
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Time", Value: fmt.Sprintf("%s - %s", ev.StartAt.Local().Format("15:04"), ev.EndAt.Local().Format("15:04")), Inline: true},
					{Name: "Location", Value: loc, Inline: true},
				},
			})
		}

		if len(embeds) > 0 {
			_, err = c.session.ChannelMessageSendComplex(gs.NotifyChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("🌅 Good morning! Here are today's events for **%s**:", sub.CalendarTitle),
				Embeds:  embeds,
			})
			if err != nil {
				slog.Error("Failed to send broadcast message", "channel_id", gs.NotifyChannelID, "error", err)
			}
		}
	}
}
