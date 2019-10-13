package store

import (
	"github.com/labstack/echo"
)

const (
	key = "store"
)

type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Store associated with this context.
func FromContext(c echo.Context) Store {
	return c.Get(key).(Store)
}

// ToContext adds the Store to this context if it supports
// the Setter interface.
func ToContext(c Setter, store Store) {
	c.Set(key, store)
}
