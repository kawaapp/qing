package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	var (
		signKey  = "....sign...key..."
		someText = "...some...text..."
	)

	token := New(SessToken, someText)
	raw, err := token.SignExpire(signKey, time.Now().Add(time.Hour*24).Unix())

	if err != nil {
		t.Error(err)
	}

	expect, err := Parse(raw, func(token *Token) (string, error) {
		return signKey, nil
	})

	if err != nil {
		t.Error(err)
	}

	if expect.Kind != token.Kind {
		t.Error("kind not equal:", expect.Kind, token.Kind)
	}
	if expect.Text != token.Text {
		t.Error("Text not equal:", expect.Text, token.Text)
	}
}

// Usage of jwt token
func TestJWTTokenCreate(t *testing.T) {
	var (
		signKey = "...sign..key..."
	)
	// create token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "ntop"
	claims["type"] = "typo"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// generate encoded token and send it as response
	tokenStr, err := token.SignedString([]byte(signKey))
	if err != nil {
		t.Error(err)
	}

	expect, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		t.Error(err)
	}

	jwtToken := expect.Claims.(jwt.MapClaims)
	if typo := jwtToken["type"].(string); typo != "typo" {
		t.Error("type not equal:", typo, "typo")
	}
	if text := jwtToken["name"].(string); text != "ntop" {
		t.Error("text not equal:", text, "ntop")
	}
}
