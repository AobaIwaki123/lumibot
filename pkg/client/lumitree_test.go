package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/AobaIwaki123/lumitree/pkg/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// mockAPIClient implements api.ClientWithResponsesInterface for testing.
type mockAPIClient struct {
	GetCalendarFunc       func(ctx context.Context, calendarId api.CalendarIdParam, reqEditors ...api.RequestEditorFn) (*api.GetCalendarResponse, error)
	GetCalendarEventsFunc func(ctx context.Context, calendarId api.CalendarIdParam, params *api.GetCalendarEventsParams, reqEditors ...api.RequestEditorFn) (*api.GetCalendarEventsResponse, error)
}

func (m *mockAPIClient) GetCalendarWithResponse(ctx context.Context, calendarId api.CalendarIdParam, reqEditors ...api.RequestEditorFn) (*api.GetCalendarResponse, error) {
	if m.GetCalendarFunc != nil {
		return m.GetCalendarFunc(ctx, calendarId, reqEditors...)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAPIClient) GetCalendarEventsWithResponse(ctx context.Context, calendarId api.CalendarIdParam, params *api.GetCalendarEventsParams, reqEditors ...api.RequestEditorFn) (*api.GetCalendarEventsResponse, error) {
	if m.GetCalendarEventsFunc != nil {
		return m.GetCalendarEventsFunc(ctx, calendarId, params, reqEditors...)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAPIClient) GetCalendarEventsICSWithResponse(ctx context.Context, calendarId api.CalendarIdParam, reqEditors ...api.RequestEditorFn) (*api.GetCalendarEventsICSResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAPIClient) GetHealthWithResponse(ctx context.Context, reqEditors ...api.RequestEditorFn) (*api.GetHealthResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func TestLumitreeClient_GetCalendar(t *testing.T) {
	tests := []struct {
		name        string
		mockResp    *api.GetCalendarResponse
		mockErr     error
		expectErr   bool
		expectTitle string
	}{
		{
			name: "success",
			mockResp: &api.GetCalendarResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200:      &api.Calendar{Title: "Test Calendar"},
			},
			expectErr:   false,
			expectTitle: "Test Calendar",
		},
		{
			name:      "network error",
			mockErr:   fmt.Errorf("network timeout"),
			expectErr: true,
		},
		{
			name: "not found",
			mockResp: &api.GetCalendarResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusNotFound},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAPIClient{
				GetCalendarFunc: func(ctx context.Context, calendarId api.CalendarIdParam, reqEditors ...api.RequestEditorFn) (*api.GetCalendarResponse, error) {
					return tt.mockResp, tt.mockErr
				},
			}
			client := &lumitreeClient{apiClient: mock}
			cal, err := client.GetCalendar(context.Background(), "test-id")

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cal == nil || cal.Title != tt.expectTitle {
					t.Errorf("expected calendar name %s, got %v", tt.expectTitle, cal.Title)
				}
			}
		})
	}
}

func TestLumitreeClient_GetEvents(t *testing.T) {
	testDate := openapi_types.Date{Time: time.Date(2026, 8, 25, 0, 0, 0, 0, time.UTC)}

	tests := []struct {
		name        string
		params      *api.GetCalendarEventsParams
		mockResp    *api.GetCalendarEventsResponse
		mockErr     error
		expectErr   bool
		expectCount int
	}{
		{
			name: "success with params",
			params: &api.GetCalendarEventsParams{
				From: &testDate,
				To:   &testDate,
			},
			mockResp: &api.GetCalendarEventsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &api.EventListResponse{
					Events: []api.Event{
						{Title: "Event 1"},
					},
				},
			},
			expectErr:   false,
			expectCount: 1,
		},
		{
			name: "success without params",
			mockResp: &api.GetCalendarEventsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &api.EventListResponse{
					Events: []api.Event{
						{Title: "Event 1"},
						{Title: "Event 2"},
					},
				},
			},
			expectErr:   false,
			expectCount: 2,
		},
		{
			name:      "network error",
			mockErr:   fmt.Errorf("network timeout"),
			expectErr: true,
		},
		{
			name: "server error",
			mockResp: &api.GetCalendarEventsResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAPIClient{
				GetCalendarEventsFunc: func(ctx context.Context, calendarId api.CalendarIdParam, params *api.GetCalendarEventsParams, reqEditors ...api.RequestEditorFn) (*api.GetCalendarEventsResponse, error) {
					if tt.params != nil {
						if params == nil || params.From != tt.params.From || params.To != tt.params.To {
							t.Errorf("params mismatch: expected %v, got %v", tt.params, params)
						}
					}
					return tt.mockResp, tt.mockErr
				},
			}
			client := &lumitreeClient{apiClient: mock}
			events, err := client.GetEvents(context.Background(), "test-id", tt.params)

			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(events) != tt.expectCount {
					t.Errorf("expected %d events, got %d", tt.expectCount, len(events))
				}
			}
		})
	}
}
