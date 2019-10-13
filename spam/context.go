package spam

import (
	"github.com/labstack/echo"
)

const (
	key = ".spamck"
)

type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Store associated with this context.
func FromContext(c echo.Context) SpamChecker {
	checker, ok := c.Get(key).(SpamChecker)
	if ok {
		return checker
	}
	return nil
}

// ToContext adds the Store to this context if it supports
// the Setter interface.
func ToContext(c Setter, checker SpamChecker) {
	c.Set(key, checker)
}
