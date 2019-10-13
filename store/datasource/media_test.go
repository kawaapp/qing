package datasource

import (
	"testing"
	"github.com/kawaapp/kawaqing/model"
)

func TestMediaCreate(t *testing.T) {
	s := beforeMedia()
	defer s.Close()

	media := model.Media {
		PostId: 1,
		Path:"[]",
	}
	if err := s.CreateMedia(&media); err != nil {
		t.Error(err)
	}
	get, err := s.GetMediaByPostId(media.PostId)
	if err != nil {
		t.Error(err)
	}
	if media.Path != get.Path {
		t.Error("media get, expect:", media.Path, " get:", get.Path)
	}

	if err := s.DeleteMediaByPostId(media.PostId); err != nil {
		t.Error(err)
	}

	_, err = s.GetMediaByPostId(media.PostId)
	if err == nil {
		t.Error("media should be deleted:", media.PostId)
	}
}

func TestGetMediaListByPostIds(t *testing.T) {
	s := beforeMedia()
	defer s.Close()

	pids := make([]int64, 100)
	for i := 0; i < 100; i++ {
		media := model.Media {
			PostId: int64(i),
			Path:"[]",
		}
		pids[i] = int64(i)
		if err := s.CreateMedia(&media); err != nil {
			t.Error(err)
		}
	}

	medias, err := s.GetMediaListByPostIds(pids)
	if err != nil {
		t.Error(err)
	}
	if len(medias) != 100 {
		t.Error("media size, expect:", 100, "get:", len(medias))
	}
}

func beforeMedia() *datasource {
	s := newTest()
	s.Exec("DELETE FROM medias")
	return s
}
