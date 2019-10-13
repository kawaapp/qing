package mwx

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/remote"
)

// Remote is a middleware function that initializes the Remote and attaches to
// the context of every http.Request.
func Remote(r remote.ClientsProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			remote.ToContext(c, r)
			return next(c)
		}
	}
}
