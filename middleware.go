package ga4m

import (
	"github.com/labstack/echo/v4"
)

// GoogleAnalyticsCookieEchoMiddleware extracts user Google Analytics
// session data from cookies and stores it in the context for later use
func GoogleAnalyticsCookieEchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(ContextKey, ParseSessionFromEchoContext(c))
			return next(c)
		}
	}
}
