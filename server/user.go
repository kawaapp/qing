package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/events"

	"net/http"
	"strconv"
	"log"
	"time"
	"fmt"
)

// 用户信息相关, 此处只提供更新和读取的方法，创建见 login.go

// Admin -

func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := store.FromContext(c).DeleteBindUser(int64(id)); err != nil {
		return fmt.Errorf("DeleteUser, %v", err)
	}
	return c.NoContent(200)
}

// 拉黑/禁言 form: 0 未定义
// blocked= 0, 1; 1拉黑, 0 取消
// silenced= 0, 1; 1禁言, 0 取消
func SetUserStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	usr, err := store.GetUser(c, int64(id))
	if err != nil {
		return c.String(404, err.Error())
	}

	if v, err := strconv.Atoi(c.QueryParam("blocked")); err == nil {
		if  v == 1 {
			usr.BlockedAt = UnixNow()
		} else {
			usr.BlockedAt = 0
		}
	}
	if v, err := strconv.Atoi(c.QueryParam("silenced")); err == nil {
		if v == 1 {
			usr.SilencedAt = UnixNow()
		} else {
			usr.SilencedAt = 0
		}
	}
	err = store.UpdateUser(c, usr)
	if err != nil {
		return fmt.Errorf("SetUserStatus, %v", err)
	}
	return c.JSON(200, usr)
}

// GetUser returns user profile by id.
func GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if u, err := store.GetUser(c, int64(id)); err != nil {
		return c.NoContent(404)
	} else {
		return jsonResp(c, 0, u)
	}
}

func Self(c echo.Context) error {
	usr := session.User(c)
	if usr == nil {
		return c.String(500, "unknown err")
	}
	return c.JSON(200, usr)
}


// UpdateUser can only update itself, it can't change others.
func UpdateUser(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	in := (map[string]interface{})(m)
	user := c.Get(".user").(*model.User)

	// Override existing value, only!
	if value, ok := in["nickname"]; ok {
		user.Nickname = value.(string)
	}
	if value, ok := in["email"]; ok {
		user.Email = value.(string)
	}
	if value, ok := in["phone"]; ok {
		user.Phone = value.(string)
	}
	if value, ok := in["avatar"]; ok {
		user.Avatar = value.(string)
	}

	if err := store.UpdateUser(c, user); err != nil {
		return fmt.Errorf("UpdateUser, %v", err)
	}
	return c.JSON(200, user)
}

// events
func usrOnUserLogin(c echo.Context, v interface{}) error  {
	usr, ok := v.(*model.User)
	if !ok {
		return typeError("User")
	}
	// Control rate, record very 5 minutes
	now := time.Now().Unix()
	if diff := now - usr.LastLogin; diff <= 5 * 60 {
		return nil
	}
	usr.LastLogin = now
	if err := store.UpdateUser(c, usr); err != nil {
		log.Println(err)
	}
	return nil
}

func init()  {
	events.Subscribe(eUserLogin, usrOnUserLogin)
}