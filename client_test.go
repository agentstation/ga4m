package ga4m

import (
	"net/http"
	"testing"
	"time"
)

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
