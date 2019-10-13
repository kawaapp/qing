package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"testing"
)

func TestCreateNotification(t *testing.T) {
	s := beforeNotification()
	defer s.Close()

	m := &model.Notification{
		EntityId: 1,
		EntityType: 0,
	}
	if err := s.CreateNotification(m); err != nil {
		t.Fatal(err)
	}

	// message
	get, err := s.GetNotificationById(m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if get.EntityId != m.EntityId {
		t.Fatal("notification create fail, expect:", m, "get:", get)
	}
}

func TestGetNotificationCount(t *testing.T) {
	s := beforeNotification()
	defer s.Close()

	m1 := model.Notification{
		EntityId: 1,
		ToId: 1,
		EntityType: int(model.NotLike),
	}
	m2 := model.Notification{
		EntityId: 2,
		ToId: 1,
		EntityType: int(model.NotComment),
	}
	m3 := model.Notification{
		EntityId: 1,
		ToId: 3,
		EntityType: int(model.NotLike),
	}

	// create 3 messages
	if err := s.CreateNotification(&m1); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateNotification(&m2); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateNotification(&m3); err != nil {
		t.Fatal(err)
	}

	// assert
	get1, err := s.GetNotificationCount(1)
	if err != nil {
		t.Fatal(err)
	}
	if get1.Comments != 1 {
		t.Fatal("get comments, expect:", 1, "get:", get1.Comments)
	}
	if get1.Favors != 2 {
		t.Fatal("get likes, expect:", 2, "get:", get1.Favors)
	}
}

func TestSetMessageAsRead(t *testing.T) {
	s := beforeNotification()
	defer s.Close()

	to := int64(1)
	m := &model.Notification{
		EntityId: 1,
	}
	if err := s.CreateNotification(m); err != nil {
		t.Error(err)
	}

	if err := s.SetNotificationReadId(0, m.ID); err != nil {
		t.Error(err)
	}

	gets, err := s.GetNotificationListType(to, 0, 0, 1000)
	if err != nil {
		t.Error(err)
	}

	for _, v := range  gets {
		if v.Status == 0 {
			t.Error("message not marked as read")
		}
	}
}

func TestSetMsgAsReadType(t *testing.T) {
	s := beforeNotification()
	defer s.Close()

	data := []model.Notification {
		{
			FromId: 1,
			ToId:   2,
			EntityType: int(model.NotLike),
		},
		{
			FromId: 1,
			ToId:   2,
			EntityType: int(model.NotComment),
		},
	}
	if err := createNotificationList(s, data); err != nil {
		t.Fatal(err)
	}

	err := s.SetNotificationReadType(2, model.NotLike)
	if err != nil {
		t.Fatal(err)
	}
	get, err := s.GetChatMessageById(data[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if get.Status != 1 {
		t.Fatal("SetNotificationReadType err, get:", get)
	}
}

func createNotificationList(s *datasource, data []model.Notification) error {
	for i := range data {
		if err := s.CreateNotification(&data[i]); err != nil {
			return err
		}
	}
	return nil
}

func beforeNotification() *datasource {
	s := newTest()
	s.Exec("DELETE FROM notifications;")
	return s
}
