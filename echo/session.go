package ega4m

import (
	"strconv"
	"strings"
	"time"

	"github.com/agentstation/ga4m"
	"github.com/labstack/echo/v4"
)

// ParseSessionFromContext returns the Google Analytics tracking data from an echo.Context
func ParseSessionFromContext(e echo.Context) ga4m.Session {
	var clientCookie, sessionCookie string

	// Get _ga cookie
	if cookie, err := e.Cookie(ga4m.ClientCookieName); err == nil {
		clientCookie = cookie.Value
	}

	// Get _ga_* session cookie
	for _, cookie := range e.Cookies() {
		if strings.HasPrefix(cookie.Name, ga4m.SessionCookieName) {
			sessionCookie = cookie.Value
			break
		}
	}

	// Parse Google Analytics cookies
	return parseGoogleAnalyticsCookies(clientCookie, sessionCookie)
}

// parseGoogleAnalyticsCookies parses Google Analytics cookies and returns the client ID, first visit timestamp, session count, and last session timestamp
func parseGoogleAnalyticsCookies(client, session string) ga4m.Session {
	var data ga4m.Session

	// Parse GA client cookie (GA1.1.71807069.1731019235)
	if parts := strings.Split(client, "."); len(parts) >= 4 {
		data.ClientID = parts[2] + "." + parts[3]
		if ts, err := strconv.ParseInt(parts[3], 10, 64); err == nil {
			data.FirstVisit = time.Unix(ts, 0)
		}
	}

	// Parse GA session cookie (GS1.1.1731019235.1.1.1731019762.0.0.0)
	if parts := strings.Split(session, "."); len(parts) >= 7 {
		if count, err := strconv.Atoi(parts[3]); err == nil {
			data.SessionCount = count
		}
		if ts, err := strconv.ParseInt(parts[5], 10, 64); err == nil {
			data.LastSession = time.Unix(ts, 0)
		}
	}

	return data
}
