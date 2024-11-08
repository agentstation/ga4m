package ga4m

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestNewClient(t *testing.T) {
	measurementID := "G-XXXXXXXXXX"
	apiSecret := "test_secret"

	client := NewClient(measurementID, apiSecret)

	if client.MeasurementID != measurementID {
		t.Errorf("Expected MeasurementID %s, got %s", measurementID, client.MeasurementID)
	}
	if client.APISecret != apiSecret {
		t.Errorf("Expected APISecret %s, got %s", apiSecret, client.APISecret)
	}
	if client.Endpoint != "https://www.google-analytics.com/mp/collect" {
		t.Errorf("Unexpected Endpoint %s", client.Endpoint)
	}
	if client.DebugEndpoint != "https://www.google-analytics.com/debug/mp/collect" {
		t.Errorf("Unexpected DebugEndpoint %s", client.DebugEndpoint)
	}
	if client.HTTPClient.(*http.Client).Timeout != 5*time.Second {
		t.Errorf("Unexpected HTTPClient timeout %s", client.HTTPClient.(*http.Client).Timeout)
	}
}

func TestSendEvent_Success(t *testing.T) {
	clientID := "123456.7654321"
	eventName := "test_event"
	params := map[string]interface{}{
		"param1": "value1",
		"param2": 2,
	}
	measurementID := "G-XXXXXXXXXX"
	apiSecret := "test_secret"

	// Create a mock HTTP client
	mockClient := &MockHTTPClient{}
	client := NewClient(measurementID, apiSecret)
	client.SetHTTPClient(mockClient)

	// Set up the mock Do function
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

		if payload.ClientID != clientID {
			t.Errorf("Expected ClientID %s, got %s", clientID, payload.ClientID)
		}
		if len(payload.Events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(payload.Events))
		}
		if payload.Events[0].Name != eventName {
			t.Errorf("Expected event name %s, got %s", eventName, payload.Events[0].Name)
		}

		// Verify that original params are included plus the engagement_time_msec
		expectedParamCount := len(params) + 1 // +1 for engagement_time_msec
		if len(payload.Events[0].Params) != expectedParamCount {
			t.Errorf("Expected %d params, got %d", expectedParamCount, len(payload.Events[0].Params))
		}

		// Return a successful response
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}

	err := client.SendEvent(clientID, eventName, params)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendEvent_InvalidEventName(t *testing.T) {
	clientID := "123456.7654321"
	eventName := "123_invalid_name" // Starts with a number, which is invalid
	params := map[string]interface{}{}

	client := NewClient("G-XXXXXXXXXX", "test_secret")

	err := client.SendEvent(clientID, eventName, params)
	if err == nil {
		t.Errorf("Expected error for invalid event name, got nil")
	}
}

func TestSendEvent_InvalidParams(t *testing.T) {
	clientID := "123456.7654321"
	eventName := "validEvent"
	params := map[string]interface{}{
		"1_invalid_param": "value", // Invalid parameter name
	}

	client := NewClient("G-XXXXXXXXXX", "test_secret")

	err := client.SendEvent(clientID, eventName, params)
	if err == nil {
		t.Errorf("Expected error for invalid parameter name, got nil")
	}
}

func TestSendEvents_Success(t *testing.T) {
	clientID := "123456.7654321"
	events := []EventParams{
		{
			Name: "event_one",
			Params: map[string]interface{}{
				"param1": "value1",
			},
		},
		{
			Name: "event_two",
			Params: map[string]interface{}{
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

		if len(payload.Events) != len(events) {
			t.Errorf("Expected %d events, got %d", len(events), len(payload.Events))
		}

		// Return a successful response
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}

	err := client.SendEvents(clientID, events)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendEvents_TooManyEvents(t *testing.T) {
	clientID := "123456.7654321"
	events := make([]EventParams, 26) // 26 events, exceeding the limit

	// Initialize events with valid event names and parameters
	for i := 0; i < 26; i++ {
		events[i] = EventParams{
			Name: fmt.Sprintf("event_%d", i),
			Params: map[string]interface{}{
				"param": fmt.Sprintf("value_%d", i),
			},
		}
	}

	client := NewClient("G-XXXXXXXXXX", "test_secret")

	err := client.SendEvents(clientID, events)
	if err == nil {
		t.Errorf("Expected error for too many events, got nil")
	}
}

func TestValidateEventName_Valid(t *testing.T) {
	eventName := "validEventName_123"

	err := validateEventName(eventName)
	if err != nil {
		t.Errorf("Expected event name to be valid, got error: %v", err)
	}
}

func TestValidateEventName_Invalid(t *testing.T) {
	eventName := "1InvalidEventName"

	err := validateEventName(eventName)
	if err == nil {
		t.Errorf("Expected event name to be invalid, got no error")
	}
}

func TestValidateParams_Valid(t *testing.T) {
	params := map[string]interface{}{
		"param_one": "value1",
		"paramTwo":  "value2",
	}

	err := validateParams(params)
	if err != nil {
		t.Errorf("Expected parameters to be valid, got error: %v", err)
	}
}

func TestValidateParams_Invalid(t *testing.T) {
	params := map[string]interface{}{
		"param-one": "value1", // Invalid character '-'
	}

	err := validateParams(params)
	if err == nil {
		t.Errorf("Expected parameters to be invalid, got no error")
	}
}

func TestParseSessionFromRequest_Valid(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  ClientCookieName,
		Value: "GA1.1.71807069.1731019235",
	})
	req.AddCookie(&http.Cookie{
		Name:  SessionCookieName + "XXXX",
		Value: "GS1.1.1731019235.1.1.1731019762.0.0.0",
	})

	session := ParseSessionFromRequest(req)

	if session.ClientID != "71807069.1731019235" {
		t.Errorf("Expected ClientID '71807069.1731019235', got '%s'", session.ClientID)
	}
	if session.SessionCount != 1 {
		t.Errorf("Expected SessionCount 1, got %d", session.SessionCount)
	}
	expectedLastSession := time.Unix(1731019762, 0)
	if !session.LastSession.Equal(expectedLastSession) {
		t.Errorf("Expected LastSession %v, got %v", expectedLastSession, session.LastSession)
	}
}

func TestParseSessionFromRequest_MissingCookies(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	session := ParseSessionFromRequest(req)

	if session.ClientID != "" {
		t.Errorf("Expected empty ClientID, got '%s'", session.ClientID)
	}
	if !session.FirstVisit.IsZero() {
		t.Errorf("Expected zero FirstVisit, got %v", session.FirstVisit)
	}
	if session.SessionCount != 0 {
		t.Errorf("Expected SessionCount 0, got %d", session.SessionCount)
	}
}

func TestSendPayload_HTTPError(t *testing.T) {
	clientID := "123456.7654321"
	eventName := "test_event"
	params := map[string]interface{}{
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

	err := client.SendEvent(clientID, eventName, params)
	if err == nil {
		t.Errorf("Expected error due to bad response, got nil")
	} else if !strings.Contains(err.Error(), "received non-OK status") {
		t.Errorf("Expected 'received non-OK status' error, got %v", err)
	}
}

func TestSendEvent_WithOptions(t *testing.T) {
	clientID := "123456.7654321"
	eventName := "test_event"
	params := map[string]interface{}{}
	userID := "user_123"
	sessionID := "session_456"
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

		if payload.UserID != userID {
			t.Errorf("Expected UserID '%s', got '%s'", userID, payload.UserID)
		}
		if payload.TimestampMicros != timestamp.UnixMicro() {
			t.Errorf("Expected TimestampMicros '%d', got '%d'", timestamp.UnixMicro(), payload.TimestampMicros)
		}
		if payload.Events[0].Params["session_id"] != sessionID {
			t.Errorf("Expected session_id '%s', got '%v'", sessionID, payload.Events[0].Params["session_id"])
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	}

	err := client.SendEvent(clientID, eventName, params,
		WithUserID(userID),
		WithSessionID(sessionID),
		WithTimestamp(timestamp),
	)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
