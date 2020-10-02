package datasource

import (
	"testing"
	"database/sql"
)

func TestLike(t *testing.T) {
	s := beforeLike()
	defer s.Close()

	//f := model.Like{
	//	Status:   1,
	//	UserID: 1,
	//	TargetTy: "post",
	//	TargetID: 3,
	//}

	if err, _ := s.CreateLike("post", 3, 1); err != nil {
		t.Fatal(err)
	}

	// assert create
	get, err := s.GetLike("post", 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	if get.Status != 1 {
		t.Fatal("CreateLike error, status != 1")
	}

	// assert delete
	if err := s.DeleteLike("post", 3, 1); err != nil {
		t.Error(err)
	}
	if l, err := s.GetLike("post", 3, 1); err != sql.ErrNoRows && l.Status != 0 {
		t.Fatal("DeleteLike error!", l)
	}
}

func TestGetLikeCount(t *testing.T) {
	// TODO
}

func TestGetLikeListUser(t *testing.T) {
	// TODO
}

func beforeLike() *datasource {
	s := newTest()
	s.Exec("DELETE FROM likes;")
	return s
}
