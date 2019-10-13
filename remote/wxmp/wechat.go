package wxmp

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/remote"
	"github.com/kawaapp/kawaqing/shared/httpx"
	"github.com/kawaapp/wsq/store"

	"encoding/json"
	"fmt"
	"log"
	"errors"
	"os"
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
)

// 注意：微信小程序和公众号的授权方式不一样，此处处理的是小程序相关的业务逻辑
// 小程序从客户端拿到 code 之后也无法取得用户信息
// 只能从客户端取得用户信息上传给服务器，服务器解密得到用户信息
// WeChat default URL:
const (
	// <AppId> <AppSecret> <code>
	defaultAPI = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	defaultURL = ""
	// <OpenId>
	getUserInfo = "https://api.weixin.qq.com/sns/userinfo?openid=%s&access_token=%s"
)

type Opts struct {
	URL    string   // GitHub server url.
	Client string   // WeChat app id
	Secret string   // WeChat secret
	Scopes []string // GitHub oauth scopes

}

func New(opts Opts) (remote.Remote, error) {
	remote := &client{
		API: defaultAPI,
		URL: defaultURL,

		Client: opts.Client,
		Secret: opts.Secret,
		Scopes: opts.Scopes,
	}
	return remote, nil
}

type client struct {
	URL      string
	Context  string
	API      string
	Client   string
	Secret   string
	Scopes   []string
	Username string
	Password string
}

// Login authenticates the session and returns the remote user details.
func (c *client) Login(ctx echo.Context) (remote.User, error) {
	// get the OAuth code // ctx.FormValue("code")
	var data = struct {
		Code string
	}{}
	ctx.Bind(&data)
	if len(data.Code) == 0 {
		// return bad request...
		return nil, errors.New("no code found")
	}

	id, secret := c.getIdSecret(ctx)
	if len(id) == 0 || len(secret) == 0 {
		return nil, errors.New("xiaochengxu id or secret is empty")
	}
	//
	resp, err := c.Exchange(id, secret, data.Code)
	if err != nil {
		return nil, err
	}
	// 返回远端数据..
	//return c.User(resp.OpenId, resp.SessionKey)
	var user RemoteUser
	user.json.OpenId = resp.OpenId
	return &user, nil
}

// Auth returns the WeChat user login for the given access token.
func (c *client) Auth() (string, error) {
	return "", nil
}


func (c *client) Decrypt(ctx echo.Context) (string, error) {
	data := new(EncryptData)
	if err := ctx.Bind(data); err != nil {
		return "", err
	}
	id, secret := c.getIdSecret(ctx)
	if len(id) == 0 || len(secret) == 0 {
		return "", errors.New("xiaochengxu id or secret is empty")
	}
	token, err := c.Exchange(id, secret, data.Code)
	if err != nil {
		return "", err
	}
	// 使用 session key 解密
	text, err := decrypt(token.SessionKey, data.Data, data.IV)
	if err != nil {
		log.Println("decrypt err, id", id, "token:", token, "data:", data)
		return "", err
	}
	return text, nil
}

func decrypt(key, data, iv string) (string, error) {
	aesKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", errors.New(fmt.Sprintf("decode session key, %s", err.Error()))
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", errors.New(fmt.Sprintf("decode iv, %s", err.Error()))
	}
	aesCipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New(fmt.Sprintf("decode data, %s", err.Error()))
	}
	if len(aesCipherText) % 16 != 0 {
		return "", errors.New("CipherText is not a multiple of the block size")
	}

	// decrypt
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(aesBlock, aesIV)
	aesPlainText := make([]byte, len(aesCipherText))
	mode.CryptBlocks(aesPlainText, aesCipherText)

	// un-padding
	aesPlainText = pkcs7UnPad(aesPlainText)
	return string(aesPlainText), nil
}

func pkcs7UnPad(b []byte) []byte {
	if b == nil || len(b) == 0 {
		return b
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return b
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return b
		}
	}
	return b[:len(b)-n]
}

func (c *client) getIdSecret(ctx echo.Context) (id, secret string) {
	if os.Getenv("STORE_LOCAL") != "" {
		id, secret = c.Client, c.Secret
	} else {
		meta, err := store.FromContext(ctx).GetMetaData()
		if err != nil {
			return
		}
		id, secret = meta["xiaocx_id"], meta["xiaocx_secret"]
	}
	return
}



// 用code从微信服务器交换 token
func (c *client) Exchange(id, secret, code string) (token AccessToken, err error) {
	urlStr := fmt.Sprintf(defaultAPI, id, secret, code)
	resp, err := httpx.Get(urlStr)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(resp, &token)
	if err != nil {
		return token, err
	}
	if len(token.OpenId) == 0 {
		return token, errors.New(string(resp))
	}
	return
}

// 微信可能没有提供这样的接口，只能在客户端上传用户信息 TODO
func (c *client) User(openid, sessionKey string) (*RemoteUser, error) {
	urlStr := fmt.Sprintf(getUserInfo, openid, sessionKey)
	resp, err := httpx.Get(urlStr)
	if err != nil {
		return nil, err
	}
	log.Println("get wechat resp:", string(resp))
	var ru RemoteUser
	if err := json.Unmarshal(resp, &ru.json); err != nil {
		return &ru, nil
	} else {
		return nil, err
	}
}

// {"session_key":"YuADihn6NytwK8Vltlb9DA==","openid":"o4GTc4mqkENdidsVaNrlv9dqpNX4"}
type AccessToken struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
}

type RemoteUser struct {
	json struct {
		OpenId    string   `json:"openid"`
		NickName  string   `json:"nickname"`
		Sex       string   `json:"sex"`
		Province  string   `json:"province"`
		City      string   `json:"city"`
		Country   string   `json:"country"`
		Avatar    string   `json:"headimgurl"`
		Privilege []string `json:"privilege"`
		UnionId   string   `json:"unionid"`
	}
}

func (ru *RemoteUser) BindId() string {
	return ru.json.OpenId
}

func (ru *RemoteUser) Name() string {
	return ru.json.NickName
}

func (ru *RemoteUser) Kind() string  {
	return "wx"
}

// Avatar
func (ru *RemoteUser) Avatar() string {
	return  ru.json.Avatar
}

// Email
func (ru *RemoteUser) Email() string {
	return ""
}

// Phone number
func (ru *RemoteUser) PhoneNumber() string {
	return ""
}

// UnionId
func (ru *RemoteUser) UnionId() string {
	return ru.json.UnionId
}

type EncryptData struct {
	Code string `json:"code"`
	Data string `json:"data"`
	IV   string `json:"iv"`
}