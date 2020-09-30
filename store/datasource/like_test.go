package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"testing"
	"math"
)

func TestLike(t *testing.T) {
	s := beforeLike()
	defer s.Close()

	f := model.Like{
		Status:   1,
		AuthorID: 1,
		PostId: 3,
	}
	if err := s.CreateLike(&f); err != nil {
		t.Error(err)
	}

	// assert create
	if get, err := s.GetLike(3, 1); err != nil || get.Status != f.Status {
		t.Error(err, "like create failed")
	}

	// assert update
	f.Status = 0
	if err := s.UpdateLike(&f); err != nil {
		t.Error(err)
	}
	if get, err := s.GetLike(3, 1); err != nil || get.Status != f.Status {
		t.Error(err, "like update failed")
	}

	// assert delete
	if err := s.DeleteLike(f.ID); err != nil {
		t.Error(err)
	}
	_, err := s.GetLike(3, 1)
	if err == nil {
		t.Error("like should be deleted.")
	}
}

func TestGetLikeCount(t *testing.T) {
	s := beforeLike()
	defer s.Close()

	var (
		entityId   = int64(3)
	)

	// lots of favor
	for i := 0; i < 100; i++ {
		f := model.Like{
			Status:     1,
			AuthorID:   int64(i),
			PostId:   entityId,
		}
		s.CreateLike(&f)
	}
	if num, err := s.GetLikeCount(entityId); err != nil || num != 100 {
		t.Error(err, "wrong num:", num)
	}
}

func TestGetLikeListUser(t *testing.T) {
	s := beforeLike()
	defer s.Close()

	// user 1 favored 3
	f1 := model.Like{
		Status:     1,
		AuthorID:   1,
		PostId:     3,
	}
	// user 1 favored 4
	f2 := model.Like{
		Status:     1,
		AuthorID:   1,
		PostId:   4,
	}
	// user 1 favored 5 and canceled
	f3 := model.Like{
		Status:     0,
		AuthorID:   1,
		PostId:     5,
	}
	s.CreateLike(&f1)
	s.CreateLike(&f2)
	s.CreateLike(&f3)

	favors, err := s.GetLikeListUser(1, 0, math.MaxInt64)
	if err != nil {
		t.Error(err)
	}
	if sz := len(favors); sz != 2 {
		t.Error("likes size, expect:", 2, "get:", sz)
	}
}

func beforeLike() *datasource {
	s := newTest()
	s.Exec("DELETE FROM likes;")
	return s
}
