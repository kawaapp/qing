package wxh5

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/shared/httpx"
	"github.com/kawaapp/kawaqing/remote"

	"fmt"
	"encoding/json"
	"net/http"
)

const (
	userInfoUrl = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	oauth2Url   = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
)

type Opts struct {
	URL    string   // GitHub server url.
	Client string   // WeChat app id
	Secret string   // WeChat secret
	Scopes []string // GitHub oauth scopes

}

func New(opts Opts) (remote.Remote, error) {
	remote := &client{
		client: opts.Client,
		secret: opts.Secret,
	}
	return remote, nil
}

// 公众号开发走的是标准的 OAuth 流程
type client struct {
	client string
	secret string
}

func (c *client) Login(ctx echo.Context) (remote.User, error) {
	in := struct {
		Code string `json:"code"`
	}{}
	if err := ctx.Bind(&in); err != nil {
		return nil, ctx.String(http.StatusBadRequest, "code not found")
	}

	token, err := c.fetchToken(in.Code)
	if err != nil {
		return nil, err
	}

	user, err := c.getUserInfo(token.AccessToken, token.OpenId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *client) getUserInfo(token, openId string) (*RemoteUser, error) {
	url := fmt.Sprintf(userInfoUrl, token, openId)
	resp, err := httpx.Get(url)
	if err != nil {
		return nil, err
	}
	data := new(RemoteUser)
	err = json.Unmarshal(resp, &data.json)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *client) fetchToken(code string) (*wxResp, error) {
	appId, secret := c.getIdSecret()
	url := fmt.Sprintf(oauth2Url, appId, secret, code)
	resp, err := httpx.Get(url)
	if err != nil {
		return nil, err
	}

	data := new(wxResp)
	err = json.Unmarshal(resp, data)
	if err != nil {
		return nil, err
	}

	if data.ErrorCode != 0 {
		return nil, fmt.Errorf("%d:%s", data.ErrorCode, data.ErrorMSg)
	}
	return data, nil
}

func (c *client) Auth() (string, error) {
	return "", nil
}

func (c *client) Decrypt(ctx echo.Context) (string, error) {
	return "", nil
}

func (c *client) getIdSecret() (string, string) {
	return c.client, c.secret
}

/*
eg:
{
  "access_token":"ACCESS_TOKEN",
  "expires_in":7200,
  "refresh_token":"REFRESH_TOKEN",
  "openid":"OPENID",
  "scope":"SCOPE"
}
 */
type wxResp struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	OpenId     string `json:"openid"`
	ExpiredIn   int64 `json:"expires_in"`
	ErrorCode   int64 `json:"errcode"`
	ErrorMSg    string `json:"errmsg"`
}

/*
eg:
{
	"openid":"obAkW5tDu2h346ryheJfBIAX-zrY",
	"nickname":"刘照云|",
	"sex":1,
	"language":"zh_CN",
	"city":"徐汇",
	"province":"上海",
	"country":"中国",
	"headimgurl":"",
    "privilege":[]
}
*/
type RemoteUser struct {
	json struct {
		OpenId    string   `json:"openid"`
		NickName  string   `json:"nickname"`
		Sex       int      `json:"-"`
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

func (*RemoteUser) Kind() string {
	return "h5"
}

func (ru *RemoteUser) Avatar() string {
	return ru.json.Avatar
}

func (*RemoteUser) Email() string {
	return ""
}

func (*RemoteUser) PhoneNumber() string {
	return ""
}

func (ru *RemoteUser) UnionId() string {
	return ru.json.UnionId
}


