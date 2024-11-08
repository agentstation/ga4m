package ga4m

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EventParams represents parameters for a GA4 event.
type EventParams struct {
	Name            string                 `json:"name"`
	Params          map[string]interface{} `json:"params,omitempty"`
	TimestampMicros int64                  `json:"timestamp_micros,omitempty"`
}

// AnalyticsEvent represents the payload structure for GA4 events.
type AnalyticsEvent struct {
	ClientID        string        `json:"client_id"`
	Events          []EventParams `json:"events"`
	UserID          string        `json:"user_id,omitempty"`
	TimestampMicros int64         `json:"timestamp_micros,omitempty"`
}

// SendEvent sends a single event to Google Analytics.
func (c *AnalyticsClient) SendEvent(clientID, eventName string, params map[string]interface{}, opts ...SendEventOption) error {
	if err := validateEventName(eventName); err != nil {
		return fmt.Errorf("invalid event name: %w", err)
	}

	if err := validateParams(params); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Apply default options.
	options := defaultSendEventOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Initialize params if nil.
	if params == nil {
		params = make(map[string]interface{})
	} else {
		// Create a copy of the params map to avoid modifying the original
		paramsCopy := make(map[string]interface{}, len(params))
		for k, v := range params {
			paramsCopy[k] = v
		}
		params = paramsCopy
	}

	// Add required session parameters if not present.
	if options.sessionID != "" {
		if _, ok := params["session_id"]; !ok {
			params["session_id"] = options.sessionID
		}
	}
	if _, ok := params["engagement_time_msec"]; !ok {
		params["engagement_time_msec"] = "100" // Default engagement time.
	}

	event := EventParams{
		Name:   eventName,
		Params: params,
	}

	if !options.timestamp.IsZero() {
		event.TimestampMicros = options.timestamp.UnixMicro()
	}

	payload := AnalyticsEvent{
		ClientID: clientID,
		Events:   []EventParams{event},
	}

	if options.userID != "" {
		payload.UserID = options.userID
	}

	if !options.timestamp.IsZero() {
		payload.TimestampMicros = options.timestamp.UnixMicro()
	}

	return c.sendPayload(payload, options)
}

// SendEvents sends multiple events in a single batch request to Google Analytics.
func (c *AnalyticsClient) SendEvents(clientID string, events []EventParams, opts ...SendEventOption) error {
	if len(events) > 25 {
		return fmt.Errorf("requests can have a maximum of 25 events")
	}

	for _, event := range events {
		if err := validateEventName(event.Name); err != nil {
			return fmt.Errorf("invalid event name '%s': %w", event.Name, err)
		}
		if err := validateParams(event.Params); err != nil {
			return fmt.Errorf("invalid parameters for event '%s': %w", event.Name, err)
		}
	}

	// Apply default options.
	options := defaultSendEventOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Add required session parameters to each event if not present.
	for i := range events {
		if events[i].Params == nil {
			events[i].Params = make(map[string]interface{})
		}
		if options.sessionID != "" {
			if _, ok := events[i].Params["session_id"]; !ok {
				events[i].Params["session_id"] = options.sessionID
			}
		}
		if _, ok := events[i].Params["engagement_time_msec"]; !ok {
			events[i].Params["engagement_time_msec"] = "100"
		}

		if !options.timestamp.IsZero() && events[i].TimestampMicros == 0 {
			events[i].TimestampMicros = options.timestamp.UnixMicro()
		}
	}

	payload := AnalyticsEvent{
		ClientID: clientID,
		Events:   events,
	}

	if options.userID != "" {
		payload.UserID = options.userID
	}

	if !options.timestamp.IsZero() {
		payload.TimestampMicros = options.timestamp.UnixMicro()
	}

	return c.sendPayload(payload, options)
}

// sendPayload handles the HTTP request to the Google Analytics endpoint.
func (c *AnalyticsClient) sendPayload(payload AnalyticsEvent, options *sendEventOptions) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	endpoint := c.Endpoint
	if options.debug {
		endpoint = c.DebugEndpoint
	}

	url := fmt.Sprintf("%s?measurement_id=%s&api_secret=%s", endpoint, c.MeasurementID, c.APISecret)

	req, err := http.NewRequestWithContext(options.ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
