package bot

import (
	"context"
	"testing"

	"github.com/AobaIwaki123/lumibot/pkg/store"
	"github.com/AobaIwaki123/lumitree/pkg/api"
)

// MockStore implements store.Store for testing.
type MockStore struct {
	AddSubscriptionFunc    func(ctx context.Context, guildID, calendarID, title string) error
	ListSubscriptionsFunc  func(ctx context.Context, guildID string) ([]store.Subscription, error)
	RemoveSubscriptionFunc func(ctx context.Context, guildID, calendarID string) error
	SetNotifyChannelFunc   func(ctx context.Context, guildID, channelID string) error
	GetGuildSettingsFunc   func(ctx context.Context, guildID string) (*store.GuildSettings, error)
	CloseFunc              func() error
}

func (m *MockStore) AddSubscription(ctx context.Context, guildID, calendarID, title string) error {
	return m.AddSubscriptionFunc(ctx, guildID, calendarID, title)
}
func (m *MockStore) ListSubscriptions(ctx context.Context, guildID string) ([]store.Subscription, error) {
	return m.ListSubscriptionsFunc(ctx, guildID)
}
func (m *MockStore) RemoveSubscription(ctx context.Context, guildID, calendarID string) error {
	return m.RemoveSubscriptionFunc(ctx, guildID, calendarID)
}
func (m *MockStore) SetNotifyChannel(ctx context.Context, guildID, channelID string) error {
	return m.SetNotifyChannelFunc(ctx, guildID, channelID)
}
func (m *MockStore) GetGuildSettings(ctx context.Context, guildID string) (*store.GuildSettings, error) {
	return m.GetGuildSettingsFunc(ctx, guildID)
}
func (m *MockStore) Close() error {
	return nil
}

// MockClient implements client.Client for testing.
type MockClient struct {
	GetCalendarFunc func(ctx context.Context, id string) (*api.Calendar, error)
	GetEventsFunc   func(ctx context.Context, id string, params *api.GetCalendarEventsParams) ([]api.Event, error)
}

func (m *MockClient) GetCalendar(ctx context.Context, id string) (*api.Calendar, error) {
	return m.GetCalendarFunc(ctx, id)
}
func (m *MockClient) GetEvents(ctx context.Context, id string, params *api.GetCalendarEventsParams) ([]api.Event, error) {
	return m.GetEventsFunc(ctx, id, params)
}

func TestNewBot(t *testing.T) {
	token := "test-token"
	mockStore := &MockStore{}
	mockClient := &MockClient{}

	b, err := New(token, mockStore, mockClient)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if b == nil {
		t.Fatal("expected bot instance, got nil")
	}

	// Token should have "Bot " prefix added automatically
	if b.Session.Token != "Bot test-token" {
		t.Errorf("expected Bot prefix in token, got %s", b.Session.Token)
	}
}
