package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"math"
	"strconv"
	"testing"
)

func TestPost(t *testing.T) {
	s := beforePost()
	defer s.Close()

	p := model.Post{
		Content: "This is the first post",
	}
	if err := s.CreatePost(&p); err != nil {
		t.Fatal(err)
	}
	if err := s.UpdatePost(&p); err != nil {
		t.Fatal(err)
	}
	getp, err := s.GetPost(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != getp.ID {
		t.Fatal("post's ID not equal")
	}
	if err := s.DeletePost(p.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetPost(p.ID); err == nil {
		t.Fatal("post should be deleted")
	}
}

func TestPostList(t *testing.T) {
	s := beforePost()
	defer s.Close()

	pid := int64(123)
	p1 := model.Post{
		DiscussionID: pid,
		Content:   "This is the first post",
	}
	p2 := model.Post{
		DiscussionID: pid,
		Content:   "This is the second post",
	}
	s.CreatePost(&p1)
	s.CreatePost(&p2)

	posts, err := s.GetPostList(pid, 0, math.MaxInt64)
	if err != nil {
		t.Error(err)
	}
	if len(posts) != 2 {
		t.Error("post's size is wrong, expect:", 2, "get:", len(posts))
	}
	if posts[1].Content != p1.Content {
		t.Error("post's name not equal")
	}
}

func TestGetPostListUser(t *testing.T) {
	s := beforePost()
	defer s.Close()

	uid := int64(1)

	p1 := model.Post{
		AuthorID:uid,
		DiscussionID: 1,
		Content:   "This is the first post",
	}
	p2 := model.Post{
		AuthorID:2,
		DiscussionID: 2,
		Content:   "This is the second post",
	}
	s.CreatePost(&p1)
	s.CreatePost(&p2)

	comments, err := s.GetPostListUser(uid, 0, math.MaxInt64)
	if err != nil {
		t.Error(err)
	}
	if sz := len(comments); sz != 1 {
		t.Error("comments size expect:", 1, "get:", sz)
	}
	if get := comments[0]; get.Content != p1.Content {
		t.Error("comment list err, expect:", p1.Content, "get:", get.Content)
	}
}

func TestPostPager(t *testing.T) {
	s := beforePost()
	defer s.Close()

	for i := 0; i < 100; i++ {
		c := model.Post{
			Content:   "This is the first post" + strconv.Itoa(i),
			DiscussionID: 1,
		}
		s.CreatePost(&c)
	}

	comments, err := s.GetPostList(1, 0, 10)
	if err != nil {
		t.Error(err)
	}
	if sz := len(comments); sz != 10 {
		t.Error("len(comments) != 10")
	}
	for i := 90; i < 100; i++ {
		expect := "This is the first post" + strconv.Itoa(i)
		if c := comments[99-i]; c.Content != expect {
			t.Error("comment paginator err", c.Content, expect)
		}
	}
}

func beforePost() *datasource {
	s := newTest()
	s.Exec("DELETE FROM posts;")
	s.Exec("DELETE FROM users;")
	s.Exec("DELETE FROM discussions;")
	return s
}
