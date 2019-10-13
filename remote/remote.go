package remote

import (
	"github.com/labstack/echo"
)

type Union interface {
	UnionId() string
}

type User interface {
	// UID
	BindId() string

	// Name
	Name() string

	// The platform of the remote agent
	Kind() string

	// Avatar
	Avatar() string

	// Email
	Email() string

	// Phone number
	PhoneNumber() string
}

type Remote interface {
	// Login authenticates the session and returns the
	// remote user details.
	Login(c echo.Context) (User, error)

	// Auth authenticates the session and returns the remote user
	// login for the given token and secret
	Auth() (string, error)

	// Decrypt
	Decrypt(c echo.Context) (string, error)
}


type ClientsProvider interface {
	SetRemote(kind string, r Remote)
	GetRemote(kind string) Remote
}

type clientsProvider struct {
	clients map[string]Remote
}

func NewProvider() ClientsProvider {
	return &clientsProvider{
		clients: make(map[string]Remote),
	}
}

func (p *clientsProvider) SetRemote(kind string, r Remote) {
	p.clients[kind] = r
}

func (p *clientsProvider) GetRemote(kind string) Remote {
	return p.clients[kind]
}

// Login authenticates the session and returns the
// remote user details.
func Login(c echo.Context, kind string) (User, error) {
	return FromContext(c).GetRemote(kind).Login(c)
}

// Auth authenticates the session and returns the remote user
// login for the given token and secret
func Auth(c echo.Context, kind string, token, secret string) (string, error) {
	return FromContext(c).GetRemote(kind).Auth()
}

// Decrypt data!
func Decrypt(c echo.Context, kind string) (string, error) {
	return FromContext(c).GetRemote(kind).Decrypt(c)
}