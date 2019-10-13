package datasource

import (
	"testing"
	"strconv"
)

func TestKV(t *testing.T)  {
	s := beforeMeta()
	defer s.Close()

	if err := s.PutKV("hello", "world"); err != nil {
		t.Error(err)
	}

	if v, err := s.GetKV("hello"); err != nil {
		t.Error("key 'hello' not found")
	} else if v != "world" {
		t.Error("value should be 'world'")
	}

	if err := s.DelKV("hello"); err != nil {
		t.Error(err)
	}

	if _, err := s.GetKV("hello"); err == nil {
		t.Error("key 'hello' should be deleted")
	}
}

func TestGetMetaData(t *testing.T) {
	s := beforeMeta()
	defer s.Close()

	kvs := map[string]interface{}{}
	for i := 0; i < 100; i++ {
		kvs[strconv.Itoa(i)] = strconv.Itoa(i)
	}
	if err := s.PutMetaData(kvs); err != nil {
		t.Error(err)
	}

	m, err := s.GetMetaData()
	if err != nil {
		t.Error(err)
	}
	if sz := len(m); sz != 100 {
		t.Error("meta table size, expect:", 100, "get:", sz)
	}
}

func TestPutIntValue(t *testing.T) {
	s := beforeMeta()
	defer s.Close()

	kvs := map[string]interface{}{}
	kvs["user_mode"] = 1
	kvs["hello"] = "hello"
	if err := s.PutMetaData(kvs); err != nil {
		t.Error(err)
	}
	m, err := s.GetMetaData()
	if err != nil {
		t.Error(err)
	}
	if v := m["user_mode"]; v != "1" {
		t.Error("user_mode expect:", 1, " get:", v)
	}
	if v := m["hello"]; v != "hello" {
		t.Error("hello expect:", "hello", " get:", v)
	}
}

func beforeMeta() *datasource {
	s := newTest()
	s.Exec("DELETE FROM metadata")
	return s
}
