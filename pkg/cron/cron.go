// Package cron provides scheduled background jobs for the bot.
package cron

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AobaIwaki123/lumibot/pkg/client"
	"github.com/AobaIwaki123/lumibot/pkg/store"
	"github.com/AobaIwaki123/lumitree/pkg/api"
	"github.com/bwmarrin/discordgo"
	openapi_types "github.com/oapi-codegen/runtime/types"
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
func New(s *discordgo.Session, st store.Store, apiClient client.Client) (*Cron, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %w", err)
	}

	c := robfigcron.New(robfigcron.WithLocation(loc))

	cronJob := &Cron{
		scheduler: c,
		session:   s,
		store:     st,
		api:       apiClient,
	}

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

	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	today := openapi_types.Date{Time: now}
	params := &api.GetCalendarEventsParams{
		From: &today,
		To:   &today,
	}

	for _, sub := range subs {
		events, err := c.api.GetEvents(ctx, sub.CalendarID, params)
		if err != nil {
			slog.Error("Failed to fetch events for cron", "calendar_id", sub.CalendarID, "error", err)
			continue
		}

		if len(events) == 0 {
			continue
		}

		var embeds []*discordgo.MessageEmbed
		for idx, ev := range events {
			if idx >= 5 {
				break
			}
			loc := "TBD"
			if ev.Location != nil && *ev.Location != "" {
				loc = *ev.Location
			}
			link := "No Link"
			if ev.Url != nil && *ev.Url != "" {
				link = *ev.Url
			}
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title: ev.Title,
				Color: 0x00b0f4,
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Date", Value: ev.StartAt.In(now.Location()).Format("2006/01/02"), Inline: true},
					{Name: "Location", Value: loc, Inline: true},
					{Name: "Link", Value: link, Inline: false},
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
