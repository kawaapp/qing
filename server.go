package main

import (
	"github.com/labstack/gommon/log"
	"github.com/kawaapp/kawaqing/remote"
	"github.com/kawaapp/kawaqing/remote/wxmp"
	"github.com/kawaapp/kawaqing/remote/wxh5"
	"github.com/kawaapp/kawaqing/router"
	"github.com/kawaapp/kawaqing/router/mwx"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/store/datasource"
	"github.com/kawaapp/kawaqing/router/mwx/spamck"
	"github.com/kawaapp/kawaqing/spam"
	"github.com/kawaapp/kawaqing/spam/triefilter"

	"io/ioutil"
	"encoding/json"
)

type AppConfig struct {
	Host struct {
		Path string `json:"path"`
		Port string `json:"port"`
	} `json:"host"`

	Store struct {
		Driver string `json:"driver"`
		Config string `json:"config"`
	} `json:"store"`

	OAuth struct {
		WeChat struct {
			Client string `json:"client"`
			Secret string `json:"secret"`
			URL    string `json:"url"`
		}
		WxH5 struct{
			Client string `json:"client"`
			Secret string `json:"secret"`
			URL    string `json:"url"`
		}
	} `json:"oauth"`

	Spam struct{
		Words string `json:"words"`
		File  string `json:"file"`
	}
}

func main() {
	appCfg, err := readConfig("app.cfg")
	if err != nil {
		log.Fatal(err)
	}

	// setup remote service
	remote, err := setupRemote(&appCfg)
	if err != nil {
		log.Fatal(err)
	}

	// setup store
	store := setupStore(&appCfg)
	
	// setup spam checker
	checker := setupSpamChecker(&appCfg)

	// start server
	e := router.Load(mwx.Store(store), mwx.Remote(remote), spamck.AttachSpamChecker(checker))
	e.Logger.Fatal(e.Start(appCfg.Host.Port))
}

// 如此数据库就可以配置起来了...
func setupStore(app *AppConfig) store.Store {
	var (
		driver = "sqlite3"
		config = ":memory:"
	)
	if app.Store.Driver != "" {
		driver = app.Store.Driver
		config = app.Store.Config
	}
	return datasource.New(
		driver,
		config,
	)
}

func setupRemote(app *AppConfig) (remote.ClientsProvider, error) {
	provider := remote.NewProvider()

	// add wx mini program client
	mp, err := setupWxMiniProgram(app)
	if err != nil {
		return nil, err
	}
	provider.SetRemote("mp", mp)

	// add wx h5 client
	h5, err := setupWxH5(app)
	if err != nil {
		return nil, err
	}
	provider.SetRemote("h5", h5)
	return provider, nil
}

func setupWxMiniProgram(app *AppConfig) (remote.Remote, error) {
	return wxmp.New(wxmp.Opts{
		URL:    app.OAuth.WeChat.URL,
		Client: app.OAuth.WeChat.Client,
		Secret: app.OAuth.WeChat.Secret,
	})
}

func setupWxH5(app *AppConfig) (remote.Remote, error) {
	return wxh5.New(wxh5.Opts{
		URL:    app.OAuth.WxH5.URL,
		Client: app.OAuth.WxH5.Client,
		Secret: app.OAuth.WxH5.Secret,
	})
}

func setupSpamChecker(app *AppConfig) (spam.SpamChecker) {
	return triefilter.New(app.Spam.Words)
}

func readConfig(file string) (app AppConfig, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &app)
	if err != nil {
		return
	}
	return
}
