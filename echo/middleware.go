package ega4m

import (
	"github.com/agentstation/ga4m"
	"github.com/labstack/echo/v4"
)

// ContextKey is the key middleware uses to store the Google Analytics session in the echo context
const ContextKey = "ga4m.session"

// GoogleAnalyticsCookieMiddleware extracts user Google Analytics
// session data from cookies and stores it in the context for later use
func GoogleAnalyticsCookieMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if session := ParseSessionFromContext(c); session != ga4m.EmptySession {
				c.Set(ContextKey, session)
			}
			return next(c)
		}
	}
}
