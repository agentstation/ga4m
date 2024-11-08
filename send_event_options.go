package ga4m

import (
	"context"
	"time"
)

// SendEventOption allows for optional parameters when sending events.
type SendEventOption func(*sendEventOptions)

type sendEventOptions struct {
	ctx       context.Context
	debug     bool
	userID    string
	timestamp time.Time
	sessionID string
}

func defaultSendEventOptions() *sendEventOptions {
	return &sendEventOptions{
		ctx:   context.Background(),
		debug: false,
	}
}

// WithContext sets a custom context for the request.
func WithContext(ctx context.Context) SendEventOption {
	return func(o *sendEventOptions) {
		o.ctx = ctx
	}
}

// WithDebug enables or disables debug mode.
func WithDebug(debug bool) SendEventOption {
	return func(o *sendEventOptions) {
		o.debug = debug
	}
}

// WithUserID sets the user ID for the event.
func WithUserID(userID string) SendEventOption {
	return func(o *sendEventOptions) {
		o.userID = userID
	}
}

// WithTimestamp sets a custom timestamp for the event.
func WithTimestamp(timestamp time.Time) SendEventOption {
	return func(o *sendEventOptions) {
		o.timestamp = timestamp
	}
}

// WithSessionID sets the session ID for the event.
func WithSessionID(sessionID string) SendEventOption {
	return func(o *sendEventOptions) {
		o.sessionID = sessionID
	}
}
