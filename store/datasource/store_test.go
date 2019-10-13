package datasource

import (
	"testing"
	"errors"
	"fmt"
	"database/sql"
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
	"log"
)

func TestTransact_Ok(t *testing.T) {
	s := newTest()
	defer s.Close()

	s.Exec("DELETE FROM users")
	s.Exec("DELETE FROM posts")

	user := &model.User{
		Login: "hello",
	}
	post := &model.Post{
		Content: "world",
	}

	defer func() {
		err := meddler.Load(s, "users", user, user.ID)
		if err != nil {
			t.Error(err)
		}
		err = meddler.Load(s, "posts", post, post.ID)
		if err != nil {
			t.Error(err)
		}
	}()

	// test panic
	s.Transact(func(tx *sql.Tx) (err error) {
		err = meddler.Insert(tx, "users", user)
		if err != nil {
			return err
		}
		err = meddler.Insert(tx, "posts", post)
		return
	})
}

func TestTransact_Panic(t *testing.T) {
	s := newTest()
	defer s.Close()

	s.Exec("DELETE FROM users")
	s.Exec("DELETE FROM posts")

	user := &model.User{
		Login: "hello",
	}
	post := &model.Post{
		Content: "world",
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("recover from panic", r)
		}
		if n, err := s.Count("SELECT * FROM users"); n != 0 || err != sql.ErrNoRows {
			t.Error("user should not exist!", err)
		}
		if n, err := s.Count("SELECT * FROM posts"); n != 0 || err != sql.ErrNoRows {
			t.Error("post should not exist!", err)
		}
	}()

	// test panic
	s.Transact(func(tx *sql.Tx) (err error) {
		err = meddler.Insert(tx, "users", user)
		if err != nil {
			return err
		}
		panic("just crash...")
		err = meddler.Insert(tx, "posts", post)
		return
	})
}

func TestTransact_Error(t *testing.T) {
	s := newTest()
	defer s.Close()

	s.Exec("DELETE FROM users")
	s.Exec("DELETE FROM posts")

	user := &model.User{
		Login: "hello",
	}
	post := &model.Post{
		Content: "world",
	}

	defer func() {
		if n, err := s.Count("SELECT * FROM users"); n != 0 || err != sql.ErrNoRows {
			t.Error("user should not exist!", err)
		}
		if n, err := s.Count("SELECT * FROM posts"); n != 0 || err != sql.ErrNoRows {
			t.Error("post should not exist!", err)
		}
	}()

	// test panic
	s.Transact(func(tx *sql.Tx) (err error) {
		err = meddler.Insert(tx, "users", user)
		if err != nil {
			return err
		}
		err = meddler.Insert(tx, "posts", post)
		if err != nil {
			return err
		}
		return errors.New("just kidding")
	})
}


func testColumns(names []string, cols []string) error {
	kv := make(map[string]bool, 0)
	for _, v := range names {
		kv[v] = true
	}
	for _, v := range cols {
		if kv[v] {
			delete(kv, v)
		}
	}
	if len(kv) != 0 {
		return errors.New(fmt.Sprintf("Miss columns: %v", kv))
	}
	return nil
}
