package ga4m

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ClientCookieName  = "_ga"
	SessionCookieName = "_ga_"
)

// Session represents the Google Analytics session tracking data for a user.
type Session struct {
	ClientID     string    // The client ID from _ga cookie.
	FirstVisit   time.Time // First visit timestamp.
	SessionCount int32     // Number of sessions.
	LastSession  time.Time // Last session timestamp.
}

// ParseSessionFromRequest parses the Google Analytics cookies from an HTTP request and returns a Session.
func ParseSessionFromRequest(r *http.Request) Session {
	var clientCookieValue, sessionCookieValue string

	// Get _ga cookie.
	if cookie, err := r.Cookie(ClientCookieName); err == nil {
		clientCookieValue = cookie.Value
	}

	// Get _ga_* session cookie.
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(cookie.Name, SessionCookieName) {
			sessionCookieValue = cookie.Value
			break
		}
	}

	// Parse Google Analytics cookies.
	return parseGoogleAnalyticsCookies(clientCookieValue, sessionCookieValue)
}

func parseGoogleAnalyticsCookies(clientCookieValue, sessionCookieValue string) Session {
	var data Session

	// Parse GA client cookie (e.g., GA1.1.71807069.1731019235).
	if parts := strings.Split(clientCookieValue, "."); len(parts) >= 4 {
		data.ClientID = parts[2] + "." + parts[3]
		if ts, err := strconv.ParseInt(parts[3], 10, 64); err == nil {
			data.FirstVisit = time.Unix(ts, 0)
		}
	}

	// Parse GA session cookie (e.g., GS1.1.1731019235.1.1.1731019762.0.0.0).
	if parts := strings.Split(sessionCookieValue, "."); len(parts) >= 7 {
		if count, err := strconv.ParseInt(parts[3], 10, 32); err == nil {
			data.SessionCount = int32(count)
		}
		if ts, err := strconv.ParseInt(parts[5], 10, 64); err == nil {
			data.LastSession = time.Unix(ts, 0)
		}
	}

	return data
}
