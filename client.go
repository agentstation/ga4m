package ga4m

import (
	"net/http"
	"time"
)

// HTTPClient interface allows for mocking of http.Client in tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AnalyticsClient is the client for sending events to Google Analytics
type AnalyticsClient struct {
	MeasurementID string
	APISecret     string
	Endpoint      string
	DebugEndpoint string
	HTTPClient    HTTPClient
}

// NewClient creates a new AnalyticsClient with the provided measurement ID and API secret
func NewClient(measurementID, apiSecret string) *AnalyticsClient {
	return &AnalyticsClient{
		MeasurementID: measurementID,
		APISecret:     apiSecret,
		Endpoint:      "https://www.google-analytics.com/mp/collect",
		DebugEndpoint: "https://www.google-analytics.com/debug/mp/collect",
		HTTPClient:    &http.Client{Timeout: 5 * time.Second},
	}
}

// SetHTTPClient allows setting a custom HTTP client
func (c *AnalyticsClient) SetHTTPClient(client HTTPClient) {
	c.HTTPClient = client
}
