package bot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/AobaIwaki123/lumitree/pkg/api"
	"github.com/bwmarrin/discordgo"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	switch data.Name {
	case "add":
		b.handleAdd(s, i, data.Options)
	case "list":
		b.handleList(s, i)
	case "remove":
		b.handleRemove(s, i, data.Options)
	case "today":
		b.handleToday(s, i)
	default:
		slog.Warn("Unknown command received", "command", data.Name)
	}
}

func (b *Bot) handleAdd(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("Failed to defer response", "error", err)
		return
	}

	calendarID := opts[0].StringValue()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cal, err := b.Client.GetCalendar(ctx, calendarID)
	if err != nil {
		b.followupErr(s, i, fmt.Sprintf("Failed to fetch calendar. Are you sure `%s` is correct?\nError: %v", calendarID, err))
		return
	}

	err = b.Store.AddSubscription(ctx, i.GuildID, calendarID, cal.Title)
	if err != nil {
		b.followupErr(s, i, fmt.Sprintf("Failed to save subscription: %v", err))
		return
	}

	_ = b.Store.SetNotifyChannel(ctx, i.GuildID, i.ChannelID)

	b.followupMsg(s, i, fmt.Sprintf("✅ Successfully subscribed to calendar **%s** (`%s`)!\nUpdates will be posted to this channel.", cal.Title, calendarID))
}

func (b *Bot) handleList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("Failed to defer response", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	subs, err := b.Store.ListSubscriptions(ctx, i.GuildID)
	if err != nil {
		b.followupErr(s, i, fmt.Sprintf("Failed to list subscriptions: %v", err))
		return
	}

	if len(subs) == 0 {
		b.followupMsg(s, i, "No calendars are currently registered in this server. Use `/add` to register one.")
		return
	}

	gs, _ := b.Store.GetGuildSettings(ctx, i.GuildID)
	channelInfo := "Not set"
	if gs != nil && gs.NotifyChannelID != "" {
		channelInfo = "<#" + gs.NotifyChannelID + ">"
	}

	var sb strings.Builder
	sb.WriteString("📋 **Registered Calendars**\n")
	for _, sub := range subs {
		sb.WriteString(fmt.Sprintf("- **%s** (`%s`)\n", sub.CalendarTitle, sub.CalendarID))
	}
	sb.WriteString(fmt.Sprintf("\n🔔 Notification Channel: %s", channelInfo))

	b.followupMsg(s, i, sb.String())
}

func (b *Bot) handleRemove(s *discordgo.Session, i *discordgo.InteractionCreate, opts []*discordgo.ApplicationCommandInteractionDataOption) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("Failed to defer response", "error", err)
		return
	}

	calendarID := opts[0].StringValue()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := b.Store.RemoveSubscription(ctx, i.GuildID, calendarID)
	if err != nil {
		b.followupErr(s, i, fmt.Sprintf("Failed to remove subscription: %v", err))
		return
	}

	b.followupMsg(s, i, fmt.Sprintf("🗑️ Successfully removed calendar `%s`.", calendarID))
}

func (b *Bot) handleToday(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		slog.Error("Failed to defer response", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	subs, err := b.Store.ListSubscriptions(ctx, i.GuildID)
	if err != nil {
		b.followupErr(s, i, "Failed to load subscriptions.")
		return
	}

	if len(subs) == 0 {
		b.followupMsg(s, i, "No calendars are currently registered in this server.")
		return
	}

	sub := subs[0]

	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	today := openapi_types.Date{Time: now}
	params := &api.GetCalendarEventsParams{
		From: &today,
		To:   &today,
	}

	events, err := b.Client.GetEvents(ctx, sub.CalendarID, params)
	if err != nil {
		b.followupErr(s, i, fmt.Sprintf("Failed to fetch events for `%s`: %v", sub.CalendarID, err))
		return
	}

	if len(events) == 0 {
		b.followupMsg(s, i, fmt.Sprintf("No events found for today in `%s`.", sub.CalendarTitle))
		return
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
				{
					Name:   "Date",
					Value:  ev.StartAt.In(now.Location()).Format("2006/01/02"),
					Inline: true,
				},
				{
					Name:   "Location",
					Value:  loc,
					Inline: true,
				},
				{
					Name:   "Link",
					Value:  link,
					Inline: false,
				},
			},
		})
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("📅 Today's Events for **%s**:", sub.CalendarTitle),
		Embeds:  embeds,
	})
	if err != nil {
		slog.Error("Failed to send followup", "error", err)
	}
}

func (b *Bot) followupMsg(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg,
	})
	if err != nil {
		slog.Error("Failed to send followup", "error", err)
	}
}

func (b *Bot) followupErr(s *discordgo.Session, i *discordgo.InteractionCreate, errMsg string) {
	b.followupMsg(s, i, "❌ "+errMsg)
}
