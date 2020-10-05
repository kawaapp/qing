package datasource

import (
	"testing"
	"github.com/kawaapp/kawaqing/model"
)

func TestCreateFavorite(t *testing.T) {
	s := beforeFavorite()
	defer s.Close()

	f := &model.Favorite{
		UserID: 1,
		DiscussionID: 2,
	}
	if err := s.CreateFavorite(f); err != nil {
		t.Fatal(err)
	}

	// get
	get, err := s.GetFavoriteUser(f.UserID, f.DiscussionID)
	if err != nil {
		t.Fatal(err)
	}
	if get.UserID != f.UserID {
		t.Fatal("GetFavorite err, expect:", f, "get:", get)
	}

	get, err = s.GetFavoriteId(f.ID)
	if err != nil {
		t.Fatal(err)
	}
	if get.UserID != f.UserID {
		t.Fatal("GetFavoriteId err, expect:", f, "get:", get)
	}

	// delete
	err = s.DeleteFavorite(f.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFavoriteListUser(t *testing.T) {
	s := beforeFavorite()
	defer s.Close()

	favorites := []model.Favorite {
		{
			UserID: 1,
			DiscussionID: 2,
		},
		{
			UserID: 1,
			DiscussionID: 3,
		},
		{
			UserID: 0,
			DiscussionID: 2,
		},
	}
	if err := createFavoriteList(s, favorites); err != nil {
		t.Error(err)
	}

	// user_id = 1, get 2 favorites
	gets, err := s.GetFavoriteListUser(1, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != 2 {
		t.Fatal("GetFavoriteListUser size err, expect:", 2, "get:", len(gets))
	}
}

func createFavoriteList(s *datasource, favorites []model.Favorite) error {
	for i := range favorites {
		if err := s.CreateFavorite(&favorites[i]); err != nil {
			return err
		}
	}
	return nil
}

func beforeFavorite() *datasource  {
	s := newTest()
	s.Exec("DELETE FROM favorites")
	return s
}