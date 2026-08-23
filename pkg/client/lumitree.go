// Package client provides the lumitree API client.
package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AobaIwaki123/lumitree/pkg/api"
)

// Client defines the interface for interacting with the lumitree API.
type Client interface {
	GetCalendar(ctx context.Context, id string) (*api.Calendar, error)
	GetEvents(ctx context.Context, id string, params *api.GetCalendarEventsParams) ([]api.Event, error)
}

type lumitreeClient struct {
	apiClient api.ClientWithResponsesInterface
}

// NewLumitreeClient creates a new lumitree API client.
func NewLumitreeClient(baseURL string) (Client, error) {
	c, err := api.NewClientWithResponses(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create lumitree api client: %w", err)
	}
	return &lumitreeClient{
		apiClient: c,
	}, nil
}

// GetCalendar fetches calendar meta-data by ID.
func (c *lumitreeClient) GetCalendar(ctx context.Context, id string) (*api.Calendar, error) {
	resp, err := c.apiClient.GetCalendarWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to call get calendar api: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}
	return resp.JSON200, nil
}

// GetEvents fetches events for a calendar ID.
func (c *lumitreeClient) GetEvents(ctx context.Context, id string, params *api.GetCalendarEventsParams) ([]api.Event, error) {
	resp, err := c.apiClient.GetCalendarEventsWithResponse(ctx, id, params)
	if err != nil {
		return nil, fmt.Errorf("failed to call get events api: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("empty response body")
	}
	return resp.JSON200.Events, nil
}
