package store_test

import (
	"context"
	"testing"

	"github.com/AobaIwaki123/lumibot/pkg/store"
)

func TestStore_Subscriptions_CRUD(t *testing.T) {
	ctx := context.Background()
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory store: %v", err)
	}
	defer func() { _ = s.Close() }()

	guildID := "guild-123"

	// 1. Initially empty
	subs, err := s.ListSubscriptions(ctx, guildID)
	if err != nil {
		t.Fatalf("unexpected error listing empty subscriptions: %v", err)
	}
	if len(subs) != 0 {
		t.Fatalf("expected 0 subscriptions, got %d", len(subs))
	}

	// 2. Add subscription
	if err := s.AddSubscription(ctx, guildID, "ilife_official", "iLiFE! Official"); err != nil {
		t.Fatalf("failed to add subscription: %v", err)
	}

	// 3. List and verify
	subs, err = s.ListSubscriptions(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to list subscriptions: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(subs))
	}
	if subs[0].CalendarID != "ilife_official" || subs[0].CalendarTitle != "iLiFE! Official" {
		t.Errorf("unexpected subscription content: %+v", subs[0])
	}

	// 4. Duplicate add updates title without error
	if err := s.AddSubscription(ctx, guildID, "ilife_official", "iLiFE! (Updated)"); err != nil {
		t.Fatalf("failed to update subscription: %v", err)
	}
	subs, _ = s.ListSubscriptions(ctx, guildID)
	if len(subs) != 1 || subs[0].CalendarTitle != "iLiFE! (Updated)" {
		t.Errorf("expected updated title, got %+v", subs[0])
	}

	// 5. Remove subscription
	if err := s.RemoveSubscription(ctx, guildID, "ilife_official"); err != nil {
		t.Fatalf("failed to remove subscription: %v", err)
	}
	subs, _ = s.ListSubscriptions(ctx, guildID)
	if len(subs) != 0 {
		t.Errorf("expected 0 subscriptions after removal, got %d", len(subs))
	}
}

func TestStore_GuildSettings(t *testing.T) {
	ctx := context.Background()
	s, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory store: %v", err)
	}
	defer func() { _ = s.Close() }()

	guildID := "guild-456"

	// Default settings for unconfigured guild
	gs, err := s.GetGuildSettings(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get default guild settings: %v", err)
	}
	if gs.NotifyChannelID != "" || gs.CronTime != "08:00" {
		t.Errorf("unexpected default settings: %+v", gs)
	}

	// Set notify channel
	if err := s.SetNotifyChannel(ctx, guildID, "channel-789"); err != nil {
		t.Fatalf("failed to set notify channel: %v", err)
	}

	gs, err = s.GetGuildSettings(ctx, guildID)
	if err != nil {
		t.Fatalf("failed to get updated settings: %v", err)
	}
	if gs.NotifyChannelID != "channel-789" {
		t.Errorf("expected notify channel 'channel-789', got '%s'", gs.NotifyChannelID)
	}
}
