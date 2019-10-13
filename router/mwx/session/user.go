package session

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/shared/token"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/events"

	"net/http"
	"database/sql"
)

const (
	userKey = ".user"
)

// User returns the currently attached user.
func User(c echo.Context) *model.User {
	v := c.Get(userKey)
	if v == nil {
		return nil
	}
	u, ok := v.(*model.User)
	if !ok {
		return nil
	}
	return u
}

// AttachUser will decode the token and get the user info from
// database, then attach it to the current context. you can use
// '.user' to get the value.
func AttachUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				user *model.User
				inner error
			)
			_, err := token.ParseRequest(c, func(t *token.Token) (string, error) {
				user, inner = store.GetUserLogin(c, t.Text)
				return user.Hash, inner
			})

			if inner != nil && inner != sql.ErrNoRows {
				return c.String(500, err.Error())
			}

			if err != nil {
				return c.String(http.StatusUnauthorized, err.Error())
			}

			if user.Blocked() {
				return c.String(http.StatusForbidden, "!user blocked")
			} else {
				c.Set(userKey, user)
				// emit event - user login
				events.Dispatch("user.login", c, user)
				return next(c)
			}
		}
	}
}
