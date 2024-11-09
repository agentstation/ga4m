package ga4m

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	// ContextKey is the key middleware uses to store the Google Analytics session in the echo context
	ContextKey = "ga4m.session"

	// google analytics cookie names
	clientCookieName  = "_ga"
	sessionCookieName = "_ga_"
)

// EmptySession is an empty Google Analytics session
var EmptySession = Session{}

// Session represents the Google Analytics session tracking data for a user.
type Session struct {
	// Client Cookie Data
	ClientID      string    // The client ID from _ga cookie.
	ClientVersion string    // The version from _ga cookie (e.g., "1")
	FirstVisit    time.Time // First visit timestamp.

	// Session Cookie Data
	SessionCount   int       // Number of sessions.
	LastSession    time.Time // Last session timestamp.
	SessionID      string    // Unique identifier for the current session
	SessionVersion string    // The version from _ga_* cookie (e.g., "1")
	IsEngaged      bool      // Indicates if the user is actively engaged
	HitCount       int       // Number of hits/interactions in the current session
	IsFirstSession bool      // Indicates if this is the user's first session
	IsNewSession   bool      // Indicates if this is a new session
}

// ParseSessionFromRequest parses the Google Analytics cookies from an HTTP request and returns a Session.
func ParseSessionFromRequest(r *http.Request) Session {
	var clientCookieValue, sessionCookieValue string

	// find _ga client cookie
	if cookie, err := r.Cookie(clientCookieName); err == nil {
		clientCookieValue = cookie.Value
	}

	// find _ga_* session cookie
	for _, cookie := range r.Cookies() {
		if strings.HasPrefix(cookie.Name, sessionCookieName) {
			sessionCookieValue = cookie.Value
			break
		}
	}

	// parse ga client and session cookies
	return parseGoogleAnalyticsCookies(clientCookieValue, sessionCookieValue)
}

// ParseSessionFromEchoContext returns the Google Analytics tracking data from an echo.Context
func ParseSessionFromEchoContext(e echo.Context) Session {
	return ParseSessionFromRequest(e.Request())
}

// parseGoogleAnalyticsCookies parses Google Analytics cookies and returns the client ID, first visit timestamp, session count, and last session timestamp
func parseGoogleAnalyticsCookies(client, session string) Session {
	var data Session

	// Handle empty inputs gracefully
	if client == "" && session == "" {
		return data
	}

	// Parse GA client cookie (GA1.1.476555468.1726969270)
	// Format: GA1.{version}.{clientID}.{timestamp}
	if client != "" {
		parts := strings.Split(client, ".")
		if len(parts) >= 4 && strings.HasPrefix(parts[0], "GA") {
			// Client Version - with bounds checking
			if len(parts) > 1 && parts[1] != "" {
				data.ClientVersion = parts[1]
			}

			// Client ID - ensure both parts exist and aren't empty
			lastIdx := len(parts) - 1
			secondLastIdx := lastIdx - 1
			if secondLastIdx >= 0 && parts[secondLastIdx] != "" &&
				lastIdx >= 0 && parts[lastIdx] != "" {
				data.ClientID = parts[secondLastIdx] + "." + parts[lastIdx]
			}

			// First Visit - validate reasonable range
			if lastIdx >= 0 {
				if ts, err := strconv.ParseInt(parts[lastIdx], 10, 64); err == nil {
					now := time.Now().Unix()
					if ts > 0 && ts <= now {
						data.FirstVisit = time.Unix(ts, 0)
					}
				}
			}
		}
	}

	// Parse GA session cookie (GS1.1.1731019235.1.1.1731019762.0.0.0)
	// Format: GS1.1.{sessionID}.{sessionCount}.{sessionEngagement}.{timestamp}.{hitCount}.{isFirst}.{isNewSession}
	if session != "" {
		parts := strings.Split(session, ".")
		if len(parts) >= 9 && strings.HasPrefix(parts[0], "GS") {
			// Session Version
			if len(parts) > 1 && parts[1] != "" {
				data.SessionVersion = parts[1]
			}

			// Session ID - ensure not empty
			if parts[2] != "" {
				data.SessionID = parts[2]
			}

			// Session count - validate non-negative
			if count, err := strconv.Atoi(parts[3]); err == nil && count >= 0 {
				data.SessionCount = count
			}

			// Session engagement (0 or 1)
			if engagement, err := strconv.Atoi(parts[4]); err == nil {
				data.IsEngaged = engagement == 1
			}

			// Timestamp of last activity - validate reasonable range
			if ts, err := strconv.ParseInt(parts[5], 10, 64); err == nil {
				now := time.Now().Unix()
				if ts > 0 && ts <= now {
					data.LastSession = time.Unix(ts, 0)
				}
			}

			// Hit count - validate non-negative
			if hits, err := strconv.Atoi(parts[6]); err == nil && hits >= 0 {
				data.HitCount = hits
			}

			// Is first session (0 or 1)
			if isFirst, err := strconv.Atoi(parts[7]); err == nil {
				data.IsFirstSession = isFirst == 1
			}

			// Is new session (0 or 1)
			if isNew, err := strconv.Atoi(parts[8]); err == nil {
				data.IsNewSession = isNew == 1
			}
		}
	}

	return data
}

// LatestSessions compares Google Analytics sessions and returns the latest one
func LatestSessions(sessions ...Session) Session {
	if len(sessions) == 0 {
		return EmptySession
	}
	latest := sessions[0]
	for _, session := range sessions[1:] {
		if session.LastSession.After(latest.LastSession) {
			latest = session
		}
	}
	return latest
}
