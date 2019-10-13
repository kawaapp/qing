package wxh5

import (
	"testing"
	"log"
	"encoding/json"
)

func TestFetchToken(t *testing.T) {
	c := &client{
		client: "wxf4642f245dad81ca",
		secret: "f19ed594e6b2f3f39a25b8bfae30d449",
	}

	resp, err := c.fetchToken("021uXCMa1Z2VKx0TjaNa1g6vMa1uXCMl")
	if err != nil {
		t.Fatal(err)
	}
	log.Println("get token:", resp)
}

func TestGetUserInfo(t *testing.T) {
	c := &client{
		client: "wxf4642f245dad81ca",
		secret: "f19ed594e6b2f3f39a25b8bfae30d449",
	}

	var (
		token  = "26_uVFhv-uKpy7D9aWz44cZsi7MmWQSU8HH5-CIJ2nyi9skgwsq80F9qbEHk4VM114BJydGXJl0zTXHC-Gy3tjQbQ"
		openid = "obAkW5tDu2h346ryheJfBIAX-zrY"
	)

	user, err := c.getUserInfo(token, openid)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("get user info:", user)
}

func TestJson(t *testing.T) {
	data := []byte(`
		{
			"sex": "1"
		}
		`)

	st := struct {
		Sex int `json:"sex,string|int"`
	}{}

	err := json.Unmarshal(data, &st)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("get:", st)
}