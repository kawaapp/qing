package datasource

import (
	"testing"
	"github.com/kawaapp/kawaqing/model"
	"time"
	"github.com/russross/meddler"
	"log"
)

func TestGetTotalUser(t *testing.T)  {
	s := beforeCounter()
	defer s.Close()

	users := []model.User{
		{
			Login: "ntop-1",
			Email: "ntop.liu@gmail.com",
		},
		{
			Login: "ntop-2",
			Email: "ntop.liu@gmail.com",
		},
	}
	if err := createUserListTime(s, users); err != nil {
		t.Error(err)
	}
	num, err := s.GetTotalUser()
	if err != nil {
		t.Error(err)
	}
	if num != len(users) {
		t.Error("getTotalUser err, expect:", len(users), "get:", num)
	}
}

func TestGetTotalDiscussion(t *testing.T) {
	s := beforeCounter()
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
	}
	if err := createDiscussionListTime(s, posts); err != nil {
		t.Error(err)
	}
	num, err := s.GetTotalDiscussion()
	if err != nil {
		t.Error(err)
	}
	if num != len(posts) {
		t.Error("getTotalPost err, expect:", len(posts), "get:", num)
	}

}

func TestGetNewUser(t *testing.T) {
	s := beforeCounter()
	defer s.Close()

	var (
		day1 = time.Now()
		day2 = time.Now().Add(-24 * time.Hour)
		day3 = time.Now().Add(-48 * time.Hour)
	)

	users := []model.User{
		{
			Login: "ntop-1",
			Email: "ntop.liu@gmail.com",
			CreatedAt: day1.Unix(),
		},
		{
			Login: "ntop-2",
			Email: "ntop.liu@gmail.com",
			CreatedAt: day2.Unix(),

		},
	}
	if err := createUserListTime(s, users); err != nil {
		t.Error(err)
	}
	num, err := s.GetNewUser(day1)
	if err != nil {
		t.Error(err)
	}
	if num != 1 {
		t.Error("getNewUser err, expect:", 1, " get:", num)
	}
	num, err = s.GetNewUser(day3)
	if err != nil {
		t.Error(err)
	}
	if num != 0 {
		t.Error("getNewUser err, expect:", 0, " get:", num)
	}

}

func TestGetNewDiscussion(t *testing.T) {
	s := beforeCounter()
	defer s.Close()

	var (
		day1 = time.Now()
		day2 = time.Now().Add(-24 * time.Hour)
		day3 = time.Now().Add(-48 * time.Hour)
	)

	posts := []model.Discussion{
		{
			Content: "Hello",
			CreatedAt: day1.Unix(),
		},
		{
			Content: "World",
			CreatedAt: day2.Unix(),
		},
	}
	if err := createDiscussionListTime(s, posts); err != nil {
		t.Error(err)
	}
	num, err := s.GetNewDiscussion(day1)
	if err != nil {
		t.Error(err)
	}
	if num != 1 {
		t.Error("getNewPost err, expect:", 1, " get:", num)
	}
	num, err = s.GetNewDiscussion(day3)
	if err != nil {
		t.Error(err)
	}
	if num != 0 {
		t.Error("getNewPost err, expect:", 0, " get:", num)
	}
}

func TestGetUserActive(t *testing.T) {
	s := beforeCounter()
	defer s.Close()

	var (
		day1 = time.Now()
		day2 = time.Now().Add(-24 * time.Hour)
		day3 = time.Now().Add(-48 * time.Hour)
	)

	users := []model.User{
		{
			Login: "ntop-1",
			Email: "ntop.liu@gmail.com",
			LastLogin: day1.Unix(),
		},
		{
			Login: "ntop-2",
			Email: "ntop.liu@gmail.com",
			LastLogin: day3.Unix(),

		},
	}
	if err := createUserListTime(s, users); err != nil {
		t.Error(err)
	}
	num, err := s.GetUserActive(day1)
	if err != nil {
		t.Error(err)
	}
	if num != 1 {
		t.Error("getUserActive err, expect:", 1, " get:", num)
	}
	num, err = s.GetUserActive(day2)
	if err != nil {
		t.Error(err)
	}
	if num != 0 {
		t.Error("getUserActive err, expect:", 0, " get:", num)
	}

}

func TestGetNewUserDaily(t *testing.T) {
	s := beforeCounter()
	defer s.Close()

	var (
		day1 = time.Now()
		day2 = time.Now().Add(-24 * time.Hour)
		day3 = time.Now().Add(-48 * time.Hour)
	)

	users := []model.User{
		{
			Login: "ntop-1",
			Email: "ntop.liu@gmail.com",
			CreatedAt: day1.Unix(),
		},
		{
			Login: "hello",
			Email: "ntop.liu@gmail.com",
			CreatedAt: day1.Add(-time.Minute).Unix(),
		},
		{
			Login: "ntop-2",
			Email: "ntop.liu@gmail.com",
			CreatedAt: day2.Unix(),
		},
	}
	if err := createUserListTime(s, users); err != nil {
		t.Error(err)
	}
	arr, err := s.GetNewUserDaily(day3, day1.Add(time.Millisecond))
	if err != nil {
		t.Error(err)
	}
	printCounterResult(arr)
	if len(arr) != 2 {
		t.Error("getNewUserDaily err, expect:", 2, " get:", len(arr))
		return
	}
	if arr[1].Count != 2 {
		t.Error("getNewUserDaily.count err, expect:", 2, "get:", arr[1].Count)
	}
}

func TestGetNewDiscussionDaily(t *testing.T) {
	s := beforeCounter()
	defer s.Close()

	var (
		day1 = time.Now()
		day2 = time.Now().Add(-24 * time.Hour)
		day3 = time.Now().Add(-48 * time.Hour)
	)

	posts := []model.Discussion{
		{
			Content: "Hello",
			CreatedAt: day1.Unix(),
		},
		{
			Content: "World",
			CreatedAt: day2.Unix(),
		},
	}
	if err := createDiscussionListTime(s, posts); err != nil {
		t.Error(err)
	}
	arr, err := s.GetNewDiscussionDaily(day3, day1.Add(time.Millisecond))
	if err != nil {
		t.Error(err)
	}
	if len(arr) != 2 {
		t.Error("getNewDiscussionDaily err, expect:", 2, " get:", len(arr))
	}
}

func TestGetActiveUserDaily(t *testing.T) {
	s := beforeCounter()
	defer s.Close()

	var (
		day1 = time.Now()
		day2 = time.Now().Add(-24 * time.Hour)
		day3 = time.Now().Add(-48 * time.Hour)
	)

	users := []model.User{
		{
			Login: "ntop-1",
			Email: "ntop.liu@gmail.com",
			LastLogin: day1.Unix(),
		},
		{
			Login: "ntop-2",
			Email: "ntop.liu@gmail.com",
			LastLogin: day2.Unix(),
		},
	}
	if err := createUserListTime(s, users); err != nil {
		t.Error(err)
	}
	arr, err := s.GetUserActiveDaily(day3, day1.Add(time.Millisecond))
	if err != nil {
		t.Error(err)
	}
	if len(arr) != 2 {
		t.Error("getActiveUserDaily err, expect:", 2, " get:", len(arr))
	}
}


func createUserListTime(s *datasource, users []model.User) error {
	for i := range users {
		if err := meddler.Insert(s, "users", &users[i]); err != nil {
			return err
		}
	}
	return nil
}

func createDiscussionListTime(s *datasource, posts []model.Discussion) error {
	for i := range posts {
		if err := meddler.Insert(s, "discussions", &posts[i]); err != nil {
			return  err
		}
	}
	return nil
}

func printCounterResult(arr []*model.DailyCount) {
	log.Println("")
	for _, v := range arr {
		log.Print(v)
	}
}


func beforeCounter() *datasource {
	s := newTest()
	s.Exec("DELETE FROM users;")
	s.Exec("DELETE FROM posts;")
	return s
}