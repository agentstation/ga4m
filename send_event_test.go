package ga4m

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of the HTTP client.
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestSendEvent_Success(t *testing.T) {
	session := Session{
		ClientID:      "123456.7654321",
		SessionID:     "session_123",
		ClientVersion: "1",
		SessionCount:  1,
		IsEngaged:     true,
	}
	eventName := "test_event"
	params := map[string]string{
		"param1": "value1",
		"param2": "2",
	}
	measurementID := "G-XXXXXXXXXX"
	apiSecret := "test_secret"

	// Create a mock HTTP client
	mockClient := &MockHTTPClient{}
	client := NewClient(measurementID, apiSecret)
	client.SetHTTPClient(mockClient)

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		// Verify the request URL
		expectedURL := "https://www.google-analytics.com/mp/collect?measurement_id=G-XXXXXXXXXX&api_secret=test_secret"
		if req.URL.String() != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
		}

		// Verify the request body
		bodyBytes, _ := io.ReadAll(req.Body)
		var payload AnalyticsEvent
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			t.Errorf("Failed to unmarshal request body: %v", err)
		}

		if payload.ClientID != session.ClientID {
			t.Errorf("Expected ClientID %s, got %s", session.ClientID, payload.ClientID)
		}
		if len(payload.Events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(payload.Events))
		}
		if payload.Events[0].Name != eventName {
			t.Errorf("Expected event name %s, got %s", eventName, payload.Events[0].Name)
		}

		// Verify that original params are included plus the engagement_time_msec and session_id
		expectedParamCount := len(params) + 2 // +2 for engagement_time_msec and session_id
		if len(payload.Events[0].Params) != expectedParamCount {
			t.Errorf("Expected %d params, got %d", expectedParamCount, len(payload.Events[0].Params))
		}

		// Verify session_id is included in params
		if payload.Events[0].Params["session_id"] != session.SessionID {
			t.Errorf("Expected session_id %s, got %s", session.SessionID, payload.Events[0].Params["session_id"])
		}

		// Optionally, verify engagement_time_msec is included
		if _, ok := payload.Events[0].Params["engagement_time_msec"]; !ok {
			t.Errorf("Expected engagement_time_msec to be included in params")
		}

		// Return a successful response
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}

	err := client.SendEvent(session, eventName, params)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendEvent_InvalidEventName(t *testing.T) {
	session := Session{
		ClientID: "123456.7654321",
	}
	eventName := "123_invalid_name" // Starts with a number, which is invalid
	params := map[string]string{}

	client := NewClient("G-XXXXXXXXXX", "test_secret")

	err := client.SendEvent(session, eventName, params)
	if err == nil {
		t.Errorf("Expected error for invalid event name, got nil")
	}
}

func TestSendEvent_InvalidParams(t *testing.T) {
	session := Session{
		ClientID: "123456.7654321",
	}
	eventName := "validEvent"
	params := map[string]string{
		"1_invalid_param": "value", // Invalid parameter name
	}

	client := NewClient("G-XXXXXXXXXX", "test_secret")

	err := client.SendEvent(session, eventName, params)
	if err == nil {
		t.Errorf("Expected error for invalid parameter name, got nil")
	}
}

func TestSendEvents_Success(t *testing.T) {
	session := Session{
		ClientID:      "123456.7654321",
		SessionID:     "session_123",
		ClientVersion: "1",
		SessionCount:  1,
		IsEngaged:     true,
	}
	events := []EventParams{
		{
			Name: "event_one",
			Params: map[string]string{
				"param1": "value1",
			},
		},
		{
			Name: "event_two",
			Params: map[string]string{
				"param2": "value2",
			},
		},
	}
	measurementID := "G-XXXXXXXXXX"
	apiSecret := "test_secret"

	mockClient := &MockHTTPClient{}
	client := NewClient(measurementID, apiSecret)
	client.SetHTTPClient(mockClient)

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		bodyBytes, _ := io.ReadAll(req.Body)
		var payload AnalyticsEvent
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			t.Errorf("Failed to unmarshal request body: %v", err)
		}

		if payload.ClientID != session.ClientID {
			t.Errorf("Expected ClientID %s, got %s", session.ClientID, payload.ClientID)
		}

		if len(payload.Events) != len(events) {
			t.Errorf("Expected %d events, got %d", len(events), len(payload.Events))
		}

		// Verify session_id is included in all events
		for _, event := range payload.Events {
			if event.Params["session_id"] != session.SessionID {
				t.Errorf("Expected session_id %s, got %s", session.SessionID, event.Params["session_id"])
			}
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}

	err := client.SendEvents(session, events)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendEvents_TooManyEvents(t *testing.T) {
	session := Session{
		ClientID:      "123456.7654321",
		SessionID:     "session_123",
		ClientVersion: "1",
	}
	events := make([]EventParams, 26) // 26 events, exceeding the limit

	// Initialize events with valid event names and parameters
	for i := 0; i < 26; i++ {
		events[i] = EventParams{
			Name: fmt.Sprintf("event_%d", i),
			Params: map[string]string{
				"param": fmt.Sprintf("value_%d", i),
			},
		}
	}

	client := NewClient("G-XXXXXXXXXX", "test_secret")

	err := client.SendEvents(session, events)
	if err == nil {
		t.Errorf("Expected error for too many events, got nil")
	}
}

func TestSendPayload_HTTPError(t *testing.T) {
	session := Session{
		ClientID:      "123456.7654321",
		SessionID:     "session_123",
		ClientVersion: "1",
	}
	eventName := "test_event"
	params := map[string]string{
		"param1": "value1",
	}

	mockClient := &MockHTTPClient{}
	client := NewClient("G-XXXXXXXXXX", "test_secret")
	client.SetHTTPClient(mockClient)

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		// Simulate a 400 Bad Request response
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader("Bad Request")),
		}, nil
	}

	err := client.SendEvent(session, eventName, params)
	if err == nil {
		t.Errorf("Expected error due to bad response, got nil")
	} else if !strings.Contains(err.Error(), "received non-OK status") {
		t.Errorf("Expected 'received non-OK status' error, got %v", err)
	}
}

func TestSendEvent_WithOptions(t *testing.T) {
	session := Session{
		ClientID:      "123456.7654321",
		SessionID:     "default_session",
		ClientVersion: "1",
	}
	eventName := "test_event"
	params := map[string]string{}
	userID := "user_123"
	sessionID := "session_456"            // This should override the session.SessionID
	timestamp := time.Unix(1609459200, 0) // Fixed timestamp: 2021-01-01 00:00:00 UTC

	mockClient := &MockHTTPClient{}
	client := NewClient("G-XXXXXXXXXX", "test_secret")
	client.SetHTTPClient(mockClient)

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		bodyBytes, _ := io.ReadAll(req.Body)
		var payload AnalyticsEvent
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			t.Errorf("Failed to unmarshal request body: %v", err)
		}

		// Verify client ID from session is used
		if payload.ClientID != session.ClientID {
			t.Errorf("Expected ClientID '%s', got '%s'", session.ClientID, payload.ClientID)
		}

		// Verify option parameters are properly set
		if payload.UserID != userID {
			t.Errorf("Expected UserID '%s', got '%s'", userID, payload.UserID)
		}
		if payload.TimestampMicros != timestamp.UnixMicro() {
			t.Errorf("Expected TimestampMicros '%d', got '%d'", timestamp.UnixMicro(), payload.TimestampMicros)
		}

		// Verify session ID from options overrides session.SessionID
		if payload.Events[0].Params["session_id"] != sessionID {
			t.Errorf("Expected session_id '%s', got '%v'", sessionID, payload.Events[0].Params["session_id"])
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}

	err := client.SendEvent(session, eventName, params,
		WithUserID(userID),
		WithSessionID(sessionID),
		WithTimestamp(timestamp),
	)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendEvent_EmptySession(t *testing.T) {
	session := Session{} // Empty session
	eventName := "test_event"
	params := map[string]string{}

	client := NewClient("G-XXXXXXXXXX", "test_secret")

	err := client.SendEvent(session, eventName, params)
	if err == nil {
		t.Error("Expected error for empty session, got nil")
	}
	if !strings.Contains(err.Error(), "must have a valid client ID") {
		t.Errorf("Expected client ID validation error, got: %v", err)
	}
}
