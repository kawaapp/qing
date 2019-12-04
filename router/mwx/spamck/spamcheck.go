package spamck

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/spam"
)


// AttachSpamChecker will find or build a SpamChecker for each app.
func AttachSpamChecker(ck spam.SpamChecker ) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			spam.ToContext(c, ck)
			return next(c)
		}
	}
}

