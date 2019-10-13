package github

//import (
//	"net/url"
//	"net"
//	"strings"
//	"yunfatie/remote"
//	"net/http"
//	"yunfatie/model"
//)
//
//const (
//	defaultURL = "https://github.com"     // Default GitHub URL
//	defaultAPI = "https://api.github.com" // Default GitHub API URL
//)
//
//// Opts defines configuration options.
//type Opts struct {
//	URL         string   // GitHub server url.
//	Context     string   // Context to display in status check
//	Client      string   // GitHub oauth client id.
//	Secret      string   // GitHub oauth client secret.
//	Scopes      []string // GitHub oauth scopes
//	Username    string   // Optional machine account username.
//	Password    string   // Optional machine account password.
//	PrivateMode bool     // GitHub is running in private mode.
//	SkipVerify  bool     // Skip ssl verification.
//	MergeRef    bool     // Clone pull requests using the merge ref.
//}
//
//// New returns a Remote implementation that integrates with a GitHub Cloud or
//// GitHub Enterprise version control hosting provider.
//func New(opts Opts) (remote.Remote, error) {
//	url, err := url.Parse(opts.URL)
//	if err != nil {
//		return nil, err
//	}
//	host, _, err := net.SplitHostPort(url.Host)
//	if err == nil {
//		url.Host = host
//	}
//	remote := &client{
//		API:         defaultAPI,
//		URL:         defaultURL,
//		Context:     opts.Context,
//		Client:      opts.Client,
//		Secret:      opts.Secret,
//		Scopes:      opts.Scopes,
//		PrivateMode: opts.PrivateMode,
//		SkipVerify:  opts.SkipVerify,
//		MergeRef:    opts.MergeRef,
//		Machine:     url.Host,
//		Username:    opts.Username,
//		Password:    opts.Password,
//	}
//	if opts.URL != defaultURL {
//		remote.URL = strings.TrimSuffix(opts.URL, "/")
//		remote.API = remote.URL + "/api/v3/"
//	}
//
//	// Hack to enable oauth2 access in older GHE
//	return remote, nil
//}
//
//type client struct {
//	URL         string
//	Context     string
//	API         string
//	Client      string
//	Secret      string
//	Scopes      []string
//	Machine     string
//	Username    string
//	Password    string
//	PrivateMode bool
//	SkipVerify  bool
//	MergeRef    bool
//}
//
//func (c *client) Login(w http.ResponseWriter, r *http.Request) (*model.User, error) {
//	panic("implement me")
//}
//
//func (c *client) Auth(token, secret string) (string, error) {
//	panic("implement me")
//}
//
//
