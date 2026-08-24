package cron

import (
	"testing"
)

func TestNewCron(t *testing.T) {
	// Simple test to ensure New doesn't panic and returns an instance.
	// We'd need extensive mocking for session, store, and client to test dailyBroadcast.
	c, err := New(nil, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Cron instance")
	}
}
