package server

import (
	"github.com/labstack/echo"
	"strconv"
	"time"
	"unicode/utf8"
	"strings"
	"github.com/kawaapp/kawaqing/model"
	"net/url"
	"encoding/json"
)

type payload struct {
	Code int 			`json:"code"`
	Message string      `json:"message,omitempty"`

	Data interface{} 	`json:"data"`
	Entities map[string]interface{} `json:"entities"`

	Total int           `json:"total,omitempty"`
	HasMore bool        `json:"hasmore,omitempty"`
}

func makePayload(code int, data interface{}) payload {
	return payload {
		Code: code, Data: data, Entities: make(map[string]interface{}),
	}
}

func errPayload(code int, err error) payload {
	return payload{
		Code: code, Message: err.Error(),
	}
}

func includes(c echo.Context, name string) bool {
	return strings.Contains(c.QueryParam("includes"), name)
}

func getPageSize(c echo.Context) (page, size int) {
	if p, err := strconv.Atoi(c.FormValue("page")); err == nil {
		page = p
	} else {
		page = 0
	}
	if sz, err := strconv.Atoi(c.FormValue("size")); sz > 0 && err == nil {
		size = sz
	} else {
		size = 20
	}
	return
}

func getQueryParams(c echo.Context) model.QueryParams {
	m := make(map[string]string)
	query := c.QueryParam("q")
	decoded, err := url.PathUnescape(query)
	if err != nil {
		return m
	}
	if err := json.Unmarshal([]byte(decoded), &m); err != nil {
		return m
	}
	return m
}

func GetQueryIntValue(value string, dft int) int {
	if v, err := strconv.Atoi(value); err == nil {
		return v
	}
	return dft
}

// bit 操作
func Clear(number, n uint64) uint64 {
	number &= ^(1 << n)
	return number
}

func Toggle(number, n uint64) uint64 {
	number ^= 1 << n
	return number
}

func Check(number, n uint64) uint64 {
	bit := (number >> n) & 1
	return bit
}

func SetBitN(number, n, x uint64) uint64 {
	number ^= (-x ^ number) & (1 << n)
	return number
}

// 时间
func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// Unix time
func UnixNow() int64 {
	return time.Now().Unix()
}

// 字符串
func MaxSubStr(str string, cut int) string {
	if len(str) < cut {
		return str
	}
	var (
		i = 0
		w = 0
		n = len(str)
	)
	for ; i < n; i = i + w {
		_, w = utf8.DecodeRuneInString(str[i:])
		if i + w > cut {
			return str[:i]
		}
	}
	return str
}

func SubStr(str string, cut int) string {
	if runes := []rune(str); len(runes) > cut {
		return string(runes[:cut])
	}
	return str
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}