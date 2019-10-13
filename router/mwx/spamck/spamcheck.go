package spamck

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/wsq/spamcheck"
)


// AttachSpamChecker will find or build a SpamChecker for each app.
func AttachSpamChecker(ck spamcheck.SpamChecker ) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			spamcheck.ToContext(c, ck)
			return next(c)
		}
	}
}

