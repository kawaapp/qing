package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/spam"
	"fmt"
	"database/sql"
)

func CreateSpamWords(c echo.Context) error {
	in := struct {
		Words string	`json:"words"`
	}{}
	if err := c.Bind(&in); err != nil {
		return c.String(400, "no words found")
	}

	// save spam words
	err := store.FromContext(c).PutMetaData(map[string]interface{}{
		"app_spamwords": in.Words,
	})
	if err != nil {
		return err
	}
	return c.NoContent(200)
}

func GetSpamWords(c echo.Context) error {
	words, err := store.GetMetaValue(c, "app_spamwords")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return c.String(200, words)
}

func SpamCheck(c echo.Context) error  {
	in := struct {
		Text string		`json:"text"`
	}{}
	if err := c.Bind(&in); err != nil {
		return c.String(400, "no text found")
	}
	ok, word := spam.FromContext(c).Validate(in.Text)
	var result string
	if ok {
		result = fmt.Sprintf(`{ "code": 0, "message": "ok"}`)
	} else {
		result = fmt.Sprintf(`{ "code": 1, "message": "fail", "word": "%s"}`, word)
	}
	return c.String(200, result)
}