package server

import (
	"testing"
	"github.com/kawaapp/wsq/model"
	"net/url"
	"log"
	"time"
)

func TestGetQueryIntValue(t *testing.T) {
	if v := GetQueryIntValue("", 123); v != 123 {
		t.Error("from '' expect:", 123, " get:", v)
	}

	if v := GetQueryIntValue("undefined", 123); v != 123 {
		t.Error("from '' expect:", 123, " get:", v)
	}

	if v := GetQueryIntValue("123", 1234); v != 123 {
		t.Error("from '' expect:", 123, " get:", v)
	}

	if v := GetQueryIntValue("-1", -1); v != -1 {
		t.Error("from '' expect:", -1, " get:", v)
	}
}

func TestBitOperation(t *testing.T) {
	var status uint64
	status = Toggle(status, 2)
	if  status != 1<<2 {
		t.Error("toggle")
	}
	status = Toggle(status, 1)
	if Check(status,2) != 1 {
		t.Error("bit-2 not 1")
	}
	if Check(status, 1) != 1 {
		t.Error("bit-1 not 1")
	}

	status = SetBitN(status, 1, 0)
	if Check(status, 1) != 0 {
		t.Error("bit-1 not 0")
	}

	status = Clear(status, 2)
	if status != 0 {
		t.Error("bit-2 not 0")
	}
}


func TestAAA(t *testing.T) {
	t.Error(model.UserSilenced)
	t.Error(model.UserBlocked)
}

func TestURLEncoder(t *testing.T) {
	v0, err := url.PathUnescape("%E5%90%88%E7%A7%9F%2F%E8%BD%AC%E7%A7%9F")
	log.Println(v0, err)
	v1, err := url.PathUnescape(v0)
	if v0 != v1 {
		t.Error("double unescapse error")
	}
}

func TestBod(t *testing.T) {
	log.Println(Bod(time.Now()))
	log.Println(time.Now())
	log.Println(Bod(time.Now()).Add(- time.Duration(2) * 24 * time.Hour))
}

func TestMaxSubStr(t *testing.T) {
	get := MaxSubStr("", 3)
	if get != "" {
		t.Error("expect: ''")
	}

	get = MaxSubStr("ABCD", 3)
	if get != "ABC" {
		t.Error("expect: abc")
	}

	get = MaxSubStr("ABCD", 5)
	if get != "ABCD" {
		t.Error("expect: abcd")
	}

	// len = 12
	get = MaxSubStr("你好世界", 0)
	if get != "" {
		t.Error("expect: ''")
	}

	get = MaxSubStr("你好世界", 5)
	if get != "你" {
		t.Error("expect: '你'")
	}

	get = MaxSubStr("你好世界", 9)
	if get != "你好世" {
		t.Error("expect: '你'")
	}
}

func TestSubStr(t *testing.T) {
	get := SubStr("你好中国", 2)
	if get != "你好" {
		t.Error("expect：你好")
	}
	get = SubStr("你好中国", 0)
	if get != "" {
		t.Error("expect：''")
	}
	get = SubStr("你好中国", 4)
	if get != "你好中国" {
		t.Error("expect：你好中国")
	}
	get = SubStr("你好中国", 5)
	if get != "你好中国" {
		t.Error("expect：你好中国")
	}

	get = SubStr("", 5)
	if get != "" {
		t.Error("expect：''")
	}
}
