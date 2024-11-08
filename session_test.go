package ga4m

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseSessionFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		clientCookie   *http.Cookie
		sessionCookie  *http.Cookie
		expectedResult Session
	}{
		{
			name: "Standard GA Cookies",
			clientCookie: &http.Cookie{
				Name:  clientCookieName,
				Value: "GA1.1.476555468.1726969270",
			},
			sessionCookie: &http.Cookie{
				Name:  sessionCookieName + "TEST",
				Value: "GS1.1.1731019235.1.1.1731019762.0.0.0",
			},
			expectedResult: Session{
				ClientID:       "476555468.1726969270",
				ClientVersion:  "1",
				FirstVisit:     time.Unix(1726969270, 0),
				SessionID:      "1731019235",
				SessionVersion: "1",
				SessionCount:   1,
				IsEngaged:      true,
				LastSession:    time.Unix(1731019762, 0),
			},
		},
		{
			name: "Extended GA Cookie Format",
			clientCookie: &http.Cookie{
				Name:  clientCookieName,
				Value: "GA1.2.1.476555468.1726969270",
			},
			sessionCookie: &http.Cookie{
				Name:  sessionCookieName + "TEST",
				Value: "GS1.1.1731019235.2.1.1731019762.0.0.0",
			},
			expectedResult: Session{
				ClientID:       "476555468.1726969270",
				ClientVersion:  "2",
				FirstVisit:     time.Unix(1726969270, 0),
				SessionID:      "1731019235",
				SessionVersion: "1",
				SessionCount:   2,
				IsEngaged:      true,
				LastSession:    time.Unix(1731019762, 0),
			},
		},
		{
			name:           "No Cookies",
			clientCookie:   nil,
			sessionCookie:  nil,
			expectedResult: EmptySession,
		},
		{
			name: "Only Client Cookie",
			clientCookie: &http.Cookie{
				Name:  clientCookieName,
				Value: "GA1.1.476555468.1726969270",
			},
			sessionCookie: nil,
			expectedResult: Session{
				ClientID:      "476555468.1726969270",
				ClientVersion: "1",
				FirstVisit:    time.Unix(1726969270, 0),
			},
		},
		{
			name:         "Only Session Cookie",
			clientCookie: nil,
			sessionCookie: &http.Cookie{
				Name:  sessionCookieName + "TEST",
				Value: "GS1.1.1731019235.1.1.1731019762.0.0.0",
			},
			expectedResult: Session{
				SessionID:      "1731019235",
				SessionVersion: "1",
				SessionCount:   1,
				IsEngaged:      true,
				LastSession:    time.Unix(1731019762, 0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.clientCookie != nil {
				req.AddCookie(tt.clientCookie)
			}
			if tt.sessionCookie != nil {
				req.AddCookie(tt.sessionCookie)
			}

			result := ParseSessionFromRequest(req)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestParseSessionFromRequest_Valid(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  clientCookieName,
		Value: "GA1.1.71807069.1731019235",
	})
	req.AddCookie(&http.Cookie{
		Name:  sessionCookieName + "XXXX",
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

func TestLatestSessions(t *testing.T) {
	now := time.Now()
	older := now.Add(-1 * time.Hour)
	oldest := now.Add(-2 * time.Hour)

	tests := []struct {
		name           string
		sessions       []Session
		expectedResult Session
	}{
		{
			name: "Multiple Sessions",
			sessions: []Session{
				{LastSession: older, ClientID: "1"},
				{LastSession: now, ClientID: "2"},
				{LastSession: oldest, ClientID: "3"},
			},
			expectedResult: Session{LastSession: now, ClientID: "2"},
		},
		{
			name:           "No Sessions",
			sessions:       []Session{},
			expectedResult: EmptySession,
		},
		{
			name: "Single Session",
			sessions: []Session{
				{LastSession: now, ClientID: "1"},
			},
			expectedResult: Session{LastSession: now, ClientID: "1"},
		},
		{
			name: "Equal Timestamps",
			sessions: []Session{
				{LastSession: now, ClientID: "1"},
				{LastSession: now, ClientID: "2"},
			},
			expectedResult: Session{LastSession: now, ClientID: "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LatestSessions(tt.sessions...)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
