package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"strconv"
	"testing"
)

func TestCreateTag(t *testing.T) {
	s := beforeTag()
	defer s.Close()

	var (
		pid  = int64(1)
		tags = []string{"apple", "orange"}
	)

	if err := s.LinkTagDiscussion(pid, tags); err != nil {
		t.Fatal(err)
	}

	for i, v := range tags {
		tag, err := s.GetTag(v)
		if err != nil {
			t.Fatal(err)
		}
		if tag.Text != tags[i] {
			t.Fatal("tag create failed, expect:", tags[i], "get:", tag.Text)
		}
	}
}

func TestCreateRelation(t *testing.T) {
	s := beforeTag()
	defer s.Close()

	var (
		pids = []int64{1, 2, 3, 4, 5}
		tags = []string{"apple", "orange"}
	)

	for _, v := range pids {
		if err := s.LinkTagDiscussion(v, tags); err != nil {
			t.Error(err)
		}
	}

	// get
	gets := make([]*model.TagDiscussion, 0)
	err := meddler.QueryAll(s, &gets, "SELECT * FROM tag_rels;")
	if err != nil {
		t.Error(err)
	}
	if expect := len(tags) * len(pids); len(gets) != expect {
		t.Error("relataion create failed, size expect:", expect, "get:", len(gets))
	}

	tag1, _ := s.GetTag("apple")
	tag2, _ := s.GetTag("orange")

	tids := []int64{tag1.ID, tag2.ID}
	// assert pid and tag
	for i, v := range gets {
		if v.DiscussionID != pids[i/2] {
			t.Error("pid not equal, expect:", pids[i/2], "get:", v.DiscussionID)
		}
		if tid := tids[i%2]; v.TagID != tid {
			t.Error("tid not equal, expect:", tid, "get:", v.TagID)
		}
	}
}

func TestGetTagList(t *testing.T) {
	s := beforeTag()
	defer s.Close()

	var (
		pid  = int64(1)
		tags = []string{"apple", "orange", "tomato", "potato"}
	)

	if err := s.LinkTagDiscussion(pid, tags); err != nil {
		t.Error(err)
	}

	gets, err := s.GetTagList()
	if err != nil {
		t.Error(err)
	}
	for len(gets) != len(tags) {
		t.Error("expect size:", len(tags), " get:", len(gets))
	}
	for i, v := range gets {
		if tags[i] != v.Text {
			t.Error("tag create fail, expect:", tags[i], "get:", v.Text)
		}
	}
}

func TestGetPostsByTag(t *testing.T) {
	s := beforeTag()
	defer s.Close()

	// create 100 posts
	posts := make([]*model.Post, 100)
	for i := 0; i < 100; i++ {
		posts[i] = &model.Post{
			Content: "hello world + " + strconv.Itoa(i),
		}
		s.CreatePost(posts[i])
	}

	// link tags
	for _, v := range posts {
		s.LinkTagDiscussion(v.ID, []string{"apple", "orange"})
	}

	apples, err := s.GetDiscussionsByTag("apple", 0, 100)
	if err != nil {
		t.Error(err)
	}
	if len(apples) != len(posts) {
		t.Error("posts size, expect:", len(posts), "get:", len(apples))
	}

	for i, v := range posts {
		get := apples[len(posts)-i-1]
		if v.Content != get.Content {
			t.Error("get posts by tag, expect:", v.Content, "get:", get.Content)
		}
	}
}

func beforeTag() *datasource{
	s := newTest()
	s.Exec("DELETE FROM tags;")
	s.Exec("DELETE FROM tag_discussions;")
	return s
}
