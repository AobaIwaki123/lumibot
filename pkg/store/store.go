package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Store provides methods to interact with the database.
type Store interface {
	AddSubscription(ctx context.Context, guildID, calendarID, title string) error
	ListSubscriptions(ctx context.Context, guildID string) ([]Subscription, error)
	RemoveSubscription(ctx context.Context, guildID, calendarID string) error
	SetNotifyChannel(ctx context.Context, guildID, channelID string) error
	GetGuildSettings(ctx context.Context, guildID string) (*GuildSettings, error)
	Close() error
}

// SQLStore is the SQLite implementation of Store.
type SQLStore struct {
	db *sql.DB
}

// Subscription represents a calendar subscription for a guild.
type Subscription struct {
	ID            int64
	GuildID       string
	CalendarID    string
	CalendarTitle string
	CreatedAt     time.Time
}

// GuildSettings represents guild-specific notification settings.
type GuildSettings struct {
	GuildID         string
	NotifyChannelID string
	AlertChannelID  string
	CronTime        string
	CreatedAt       time.Time
}

// New creates and initializes a Store with the given database path.
func New(dbPath string) (Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	s := &SQLStore{db: db}
	if err := s.initSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return s, nil
}

// Close closes the underlying database connection.
func (s *SQLStore) Close() error {
	return s.db.Close()
}

func (s *SQLStore) initSchema(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS guild_settings (
		guild_id TEXT PRIMARY KEY,
		notify_channel_id TEXT NOT NULL DEFAULT '',
		alert_channel_id TEXT NOT NULL DEFAULT '',
		cron_time TEXT NOT NULL DEFAULT '08:00',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS subscriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guild_id TEXT NOT NULL,
		calendar_id TEXT NOT NULL,
		calendar_title TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(guild_id, calendar_id)
	);

	CREATE INDEX IF NOT EXISTS idx_subscriptions_guild_id ON subscriptions(guild_id);
	`
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SQLStore) AddSubscription(ctx context.Context, guildID, calendarID, title string) error {
	query := `
	INSERT INTO subscriptions (guild_id, calendar_id, calendar_title)
	VALUES (?, ?, ?)
	ON CONFLICT(guild_id, calendar_id) DO UPDATE SET calendar_title=excluded.calendar_title;
	`
	_, err := s.db.ExecContext(ctx, query, guildID, calendarID, title)
	if err != nil {
		return fmt.Errorf("failed to add subscription: %w", err)
	}
	return nil
}

func (s *SQLStore) ListSubscriptions(ctx context.Context, guildID string) ([]Subscription, error) {
	query := `
	SELECT id, guild_id, calendar_id, calendar_title, created_at
	FROM subscriptions
	WHERE guild_id = ?
	ORDER BY id ASC;
	`
	rows, err := s.db.QueryContext(ctx, query, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []Subscription
	for rows.Next() {
		var sub Subscription
		if err := rows.Scan(&sub.ID, &sub.GuildID, &sub.CalendarID, &sub.CalendarTitle, &sub.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		results = append(results, sub)
	}
	return results, rows.Err()
}

func (s *SQLStore) RemoveSubscription(ctx context.Context, guildID, calendarID string) error {
	query := `DELETE FROM subscriptions WHERE guild_id = ? AND calendar_id = ?;`
	_, err := s.db.ExecContext(ctx, query, guildID, calendarID)
	if err != nil {
		return fmt.Errorf("failed to remove subscription: %w", err)
	}
	return nil
}

func (s *SQLStore) SetNotifyChannel(ctx context.Context, guildID, channelID string) error {
	query := `
	INSERT INTO guild_settings (guild_id, notify_channel_id)
	VALUES (?, ?)
	ON CONFLICT(guild_id) DO UPDATE SET notify_channel_id=excluded.notify_channel_id;
	`
	_, err := s.db.ExecContext(ctx, query, guildID, channelID)
	if err != nil {
		return fmt.Errorf("failed to set notify channel: %w", err)
	}
	return nil
}

func (s *SQLStore) GetGuildSettings(ctx context.Context, guildID string) (*GuildSettings, error) {
	query := `
	SELECT guild_id, notify_channel_id, alert_channel_id, cron_time, created_at
	FROM guild_settings
	WHERE guild_id = ?;
	`
	var gs GuildSettings
	err := s.db.QueryRowContext(ctx, query, guildID).Scan(
		&gs.GuildID, &gs.NotifyChannelID, &gs.AlertChannelID, &gs.CronTime, &gs.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return &GuildSettings{GuildID: guildID, CronTime: "08:00"}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get guild settings: %w", err)
	}
	return &gs, nil
}
