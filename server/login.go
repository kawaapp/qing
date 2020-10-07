package server

import (
	"github.com/gorilla/securecookie"
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/remote"
	"github.com/kawaapp/kawaqing/shared/auth"
	"github.com/kawaapp/kawaqing/shared/token"
	"github.com/kawaapp/kawaqing/store"

	"encoding/base32"
	"errors"
	"net/http"
	"time"
	"database/sql"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"fmt"
)

type loginUser struct {
	Code     string `json:"code"`
	Phone    string `json:"phone"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

// 此处实现了4种注册/登录方式
// 1. 用户名/密码注册登录
// 2. 手机号/密码注册登录
// 3. 小程序授权登录
// 4. 公众号H5授权登录
// 注册登录需要先注册后登录，授权登录直接登录（如果没有账号自动创建）

// Register creates a new user in database.
func Register(c echo.Context) error {
	in := &loginUser{}
	if err := c.Bind(in); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user := &model.User{
		Login: in.Name,
		Hash: base32.StdEncoding.EncodeToString(
			securecookie.GenerateRandomKey(32),
		),
		PasswordHash: auth.HashPassword(in.Password),
	}

	if err := user.Validate(); err != nil {
		return c.String(http.StatusBadRequest, err.Error() + user.Login)
	}

	if err := store.CreateUser(c, user); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

// HandleLogin validates the user info, if success return an access token.
func HandleLogin(c echo.Context) error {
	in := &loginUser{}
	if err := c.Bind(in); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var (
		user *model.User
		err error
	)
	switch c.QueryParam("by") {
	case "sms":
		user, err = store.GetUserPhone(c, in.Phone)
	default:
		user, err = store.GetUserLogin(c, in.Name)
	}

	// validate user login
	if err == sql.ErrNoRows {
		return c.String(http.StatusBadRequest, LoginError.Error())
	}
	if err != nil {
		return err
	}
	if err := auth.CheckPassword(in.Password, user.PasswordHash); err != nil {
		return c.String(http.StatusBadRequest, LoginError.Error())
	}

	// if user is blocked
	if user.Blocked() {
		return c.String(http.StatusForbidden, "!user blocked")
	}

	config := getConfig(c)

	// generate token
	exp := time.Now().Add(config.Server.SessionExpires).Unix()
	token := token.New(token.SessToken, user.Login)
	raw, err := token.SignExpire(user.Hash, exp)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tokenPayload{
		Access:  raw,
		Expires: exp - time.Now().Unix(),
	})
}

// HandleAuth will create a new user, if no user found.
func HandleAuthMp(c echo.Context) error {
	user, err := remote.Login(c, "mp")
	if err != nil {
		return err
	}
	return handleAuthLogin(c, user)
}

func HandleAuthH5(c echo.Context) error {
	user, err := remote.Login(c, "h5")
	if err != nil {
		return err
	}
	return handleAuthLogin(c, user)
}

func handleAuthLogin(c echo.Context, ru remote.User) error {
	var (
		bindid = ru.BindId()
		unionid string
	)
	if union, ok := ru.(remote.Union); ok {
		unionid = union.UnionId()
	}
	user, err := tryGetUser(c, bindid, unionid)

	// create new user if not exist
	if err == sql.ErrNoRows {
		// create user
		user = &model.User {
			Login: GenerateUserId(ru.Kind()),
			Hash: base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32),
			),
			Nickname: ru.Name(),
			Avatar: ru.Avatar(),
			Email: ru.Email(),
			Phone: ru.PhoneNumber(),
		}
		bind := &model.UserBind{
			Kind: ru.Kind(),
			BindId: bindid,
			UnionId: unionid,
		}
		err = store.FromContext(c).CreateBindUser(user, bind)
	}

	// other internal error
	if err != nil {
		return err
	}

	// if user is blocked
	if user.Blocked() {
		return c.String(http.StatusForbidden, "!user blocked")
	}
	config := getConfig(c)

	// generate token
	exp := time.Now().Add(config.Server.SessionExpires).Unix()
	token := token.New(token.SessToken, user.Login)
	raw, err := token.SignExpire(user.Hash, exp)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tokenPayload{
		Access:  raw,
		Expires: exp - time.Now().Unix(),
	})
}

func tryGetUser(c echo.Context, openid, unionid string) (*model.User, error) {
	user, err := store.GetUserBind(c, openid)
	if err == nil {
		return user, nil
	}
	if len(unionid) > 0 {
		user, err = store.GetUserUnion(c, unionid)
	}
	return user, err
}

func HandleLogout(c echo.Context) error {
	return nil
}

func HandlePasswordReset(c echo.Context) error {
	in := struct {
		PasswordOld string `json:"password_old"`
		PasswordNew string `json:"password_new"`
	}{}
	if err := c.Bind(&in); err != nil {
		return c.String(200, err.Error())
	}
	usr := session.User(c)
	if err := auth.CheckPassword(in.PasswordOld, usr.PasswordHash); err != nil {
		return errorResp(c, 4001, fmt.Errorf("password error"))
	}
	hash := auth.HashPassword(in.PasswordNew)
	usr.PasswordHash = hash
	if err := store.UpdateUser(c, usr); err != nil {
		return fmt.Errorf("HandleReset, %v", err)
	}
	return c.NoContent(200)
}

func GenerateUserId(prefix string) (string) {
	uid := base32.StdEncoding.EncodeToString(
		securecookie.GenerateRandomKey(8),
	)
	return prefix+uid
}

func getConfig(c echo.Context) *model.Settings {
	return &model.Settings{
		Server: struct{ SessionExpires time.Duration }{SessionExpires: 90 * 24 * time.Hour},
	}
}

type tokenPayload struct {
	Access  string `json:"access_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	Expires int64  `json:"expires_in,omitempty"`
}

var LoginError = errors.New("username or password is wrong")
