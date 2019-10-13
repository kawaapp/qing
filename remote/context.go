package remote

import "github.com/labstack/echo"

const key = "remote"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Remote associated with this context.
func FromContext(c echo.Context) ClientsProvider {
	return c.Get(key).(ClientsProvider)
}

// ToContext adds the Remote to this context if it supports
// the Setter interface.
func ToContext(c Setter, r ClientsProvider) {
	c.Set(key, r)
}
