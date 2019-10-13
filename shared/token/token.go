package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type SecretFunc func(*Token) (string, error)

const (
	UserToken  = "user"
	SessToken  = "sess"
	HookToken  = "hook"
	CsrfToken  = "csrf"
	AgentToken = "agent"
)

// Default algorithm to sign JWT tokens.
const (
	SignerAlgo = "HS256"
)

type Token struct {
	Kind string
	Text string
}

func Parse(raw string, fn SecretFunc) (*Token, error) {
	token := &Token{}
	parsed, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
		// validate the correct algorithm is being used
		if t.Method.Alg() != SignerAlgo {
			return nil, jwt.ErrSignatureInvalid
		}

		// extract the token kind and cast to
		// the expected type
		kindV, ok := t.Claims.(jwt.MapClaims)["type"]
		if !ok {
			return nil, jwt.ValidationError{}
		}
		token.Kind, _ = kindV.(string)

		// extract the token value and cast to
		// expected type.
		textV, ok := t.Claims.(jwt.MapClaims)["text"]
		if !ok {
			return nil, jwt.ValidationError{}
		}
		token.Text, _ = textV.(string)

		// invoke the callback func to retrieve
		// the secret key used to verify
		signKey, err := fn(token)
		return []byte(signKey), err
	})

	if err != nil {
		return nil, err
	} else if !parsed.Valid {
		return nil, jwt.ValidationError{}
	}
	return token, nil
}

// ParseRequest will extract the token from echo.Context
// Note: it need echo's Token-Middleware enabled.
func ParseRequest(c echo.Context, fn SecretFunc) (*Token, error) {
	var bearer = c.Request().Header.Get("Authorization")
	var token string

	// first we attempt to get the token from the
	// authorization header.
	if len(bearer) != 0 {
		fmt.Sscanf(bearer, "Bearer %s", &token)
		return Parse(token, fn)
	}

	// then we attempt to get the token from the
	// access_token url query parameter
	token = c.FormValue("access_token")
	if len(token) != 0 {
		return Parse(token, fn)
	}

	// and finally we attempt to get the token from
	// the user session cookie
	cookie, err := c.Cookie("user_sess")
	if err != nil {
		return nil, err
	}
	return Parse(cookie.Value, fn)
}

func New(kind, text string) *Token {
	return &Token{Kind: kind, Text: text}
}

// Sign signs the token using the given secret hash
// and returns the string value.
func (t *Token) Sign(secret string) (string, error) {
	return t.SignExpire(secret, 0)
}

// Sign signs the token string using the given secret hash
// with an expiration date.
func (t *Token) SignExpire(secret string, exp int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["type"] = t.Kind
	claims["text"] = t.Text
	if exp > 0 {
		claims["exp"] = float64(exp)
	}
	return token.SignedString([]byte(secret))
}
