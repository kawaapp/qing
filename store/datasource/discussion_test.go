package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"math"
	"strconv"
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
	posts, err := s.GetDiscussionList(0, math.MaxInt64, "")
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

func TestGetDiscussionListCommentCount(t *testing.T) {
	s := beforeDiscussion()
	defer s.Close()

	posts := []model.Discussion{
		{
			Content: "Hello",
			CommentCount: 2,
		},
		{
			Content: "World",
			CommentCount: 1,
		},
		{
			Content: "Foo",
			CommentCount: 3,
		},
		{
			Content: "bar",
			CommentCount: 1,
		},
	}
	if err := createDiscussionList(s, posts); err != nil {
		t.Error(err)
	}

	// First page
	get, err := s.GetDiscussionListCommentCount(0, 2)
	if err != nil {
		t.Error(err)
	}
	if len(get) != 2 {
		t.Error("1. get post by comment-count, size get:", len(get))
		return
	}
	if get[0].Content != posts[2].Content || get[1].Content != posts[0].Content {
		t.Error("1.get posts by comment-count, get:", get[0], get[1])
	}

	// Next page
	get, err = s.GetDiscussionListCommentCount(1, 2)
	if len(get) != 2 {
		t.Error("1. get post by comment-count, size get:", len(get))
		return
	}
	if get[0].Content != posts[3].Content || get[1].Content != posts[1].Content {
		t.Error("2. get posts by comment-count, get:", get[0], get[1])
	}
}


func TestGetDiscussionListUser(t *testing.T) {
	s := beforeDiscussion()
	defer s.Close()

	uid := int64(1)
	p1 := model.Discussion{
		AuthorID:uid,
		Title:   "hello world",
		Content: "This is the first post",
	}
	p2 := model.Discussion{
		AuthorID:uid+1,
		Title:   "hello again",
		Content: "This is the second post",
	}
	s.CreateDiscussion(&p1)
	s.CreateDiscussion(&p2)

	// 降序
	posts, err := s.GetDiscussionListUser(uid, 0, math.MaxInt64)
	if err != nil {
		t.Error(err)
	}
	if len(posts) != 1 {
		t.Error("post's size is wrong")
	}
	if p := posts[0]; p.Title != p1.Title {
		t.Error("list post err, expect:", p1.Title, "get:", p.Title)
	}
}

// pagination
func TestDiscussionPager(t *testing.T) {
	s := beforeDiscussion()
	defer s.Close()

	for i := 0; i < 100; i++ {
		p := model.Discussion{
			Title:   "hello world" + strconv.Itoa(i),
			Content: "This is the first post",
		}
		s.CreateDiscussion(&p)
	}

	posts, err := s.GetDiscussionList(0, 10, "")
	if err != nil {
		t.Error(err)
	}
	if len(posts) != 10 {
		t.Error("len(post) != 10 ")
	}
	for i := 90; i < 100; i++ {
		expect := "hello world" + strconv.Itoa(i)
		if p := posts[99-i]; p.Title != expect {
			t.Error("post pagination err", p.Title, expect)
		}
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
