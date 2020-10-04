package datasource

import (
	"testing"
	"github.com/kawaapp/kawaqing/model"
)

func TestCreateFollow(t *testing.T) {
	s := beforeFollow()
	defer s.Close()

	f := &model.Follow{
		UserId: 1,
		FollowerId: 2,
	}
	if err := s.CreateFollow(f); err != nil {
		t.Fatal(err)
	}

	if _, err := s.GetFollow(1, 2); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetFollow(2, 1); err == nil {
		t.Fatal("not following between user 1 and 2")
	}
}

func TestGetFollowRelation(t *testing.T) {
	s := beforeFollow()
	defer s.Close()

	users := []model.User {
		{
			Login: "hello",
		},
		{
			Login: "world",
		},
		{
			Login: "foobar",
		},
	}
	if err := createUserList(s, users); err != nil {
		t.Fatal(err)
	}

	// create follower, 1,2 follow 0, 1 follow 2
	s.CreateFollow(&model.Follow{
		UserId: users[0].ID, FollowerId: users[1].ID,
	})
	s.CreateFollow(&model.Follow{
		UserId: users[0].ID, FollowerId: users[2].ID,
	})
	s.CreateFollow(&model.Follow{
		UserId: users[2].ID, FollowerId: users[1].ID,
	})

	// assert follower
	follower, err := s.GetFollowerList(users[0].ID, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(follower) != 2 {
		t.Fatal("GetFollowerList size err, expect:", 2, "get:", len(follower))
	}

	// 1 follow 0, 2
	following, err := s.GetFollowingList(users[1].ID, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(following) != 2 {
		t.Fatal("GetFollowingList size err, expect:", 2, "get:", len(following))
	}

	// 2 follow 0
	following, err = s.GetFollowingList(users[2].ID, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(following) != 1 {
		t.Fatal("GetFollowingList size err, expect:", 1, "get:", len(following))
	}
}

func beforeFollow() *datasource  {
	s := newTest()
	s.Exec("DELETE FROM users")
	s.Exec("DELETE FROM follows")
	return s
}