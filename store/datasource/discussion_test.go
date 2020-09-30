package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"testing"
)

func TestDiscussion(t *testing.T) {
	s := beforeDiscussion()
	defer s.Close()

	p := model.Discussion{
		Title:   "hello world",
		Content: "This is the first post",
	}
	if err := s.CreateDiscussion(&p); err != nil {
		t.Error(err)
	}
	if err := s.UpdateDiscussion(&p); err != nil {
		t.Error(err)
	}
	getp, err := s.GetDiscussion(p.ID)
	if err != nil {
		t.Error(err)
	}
	if p.ID != getp.ID {
		t.Error("post's ID not equal")
	}
	if err := s.DeleteDiscussion(p.ID); err != nil {
		t.Error(err)
	}
	if _, err := s.GetDiscussion(p.ID); err == nil {
		t.Error("post should be deleted")
	}
}

func TestDiscussionList(t *testing.T) {
	s := beforeDiscussion()
	defer s.Close()

	p1 := model.Discussion{
		Title:   "hello world",
		Content: "This is the first post",
	}
	p2 := model.Discussion{
		Title:   "hello again",
		Content: "This is the second post",
	}
	s.CreateDiscussion(&p1)
	s.CreateDiscussion(&p2)

	// 降序
	posts, err := s.GetDiscussionList(nil, 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 {
		t.Fatal("post's size is wrong, get:", len(posts))
	}
	if posts[1].Title != p1.Title {
		t.Error("post's name not equal")
	}
}

func beforeDiscussion() *datasource {
	s := newTest()
	s.Exec("DELETE FROM users;")
	s.Exec("DELETE FROM posts;")
	s.Exec("DELETE FROM favors;")
	s.Exec("DELETE FROM comments")
	return s
}
