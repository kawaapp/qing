package mwx

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/store"
)

// Store is a middleware function that initializes the Data store and attaches to
// the context of every http.Request.
func Store(v store.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			store.ToContext(c, v)
			return next(c)
		}
	}
}
