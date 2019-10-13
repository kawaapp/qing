package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"testing"
)

func TestChatMessageCreate(t *testing.T) {
	s := beforeChat()
	defer s.Close()

	m := &model.Chat {
		FromId: 1,
		ToId: 2,
		Content: "Hello",
	}

	if err := s.CreateChatMessage(m); err != nil {
		t.Fatal(err)
	}

	get, err := s.GetChatMessageById(m.ID)
	if err != nil {
		t.Fatal(err)
	}
	if get.FromId != m.FromId {
		t.Fatal("get chat message err, get:", get, "expect:", m)
	}
}

func TestSetChatMsgAsRead(t *testing.T) {
	s := beforeChat()
	defer s.Close()

	data := []model.Chat {
		{
			FromId: 1,
			ToId:   2,
			Content: "Hello",
		},
		{
			FromId: 1,
			ToId:   2,
			Content: "Hello",
		},
		{
			FromId: 1,
			ToId:   3,
			Content: "Hello",
		},
	}
	if err := createChatMessageList(s, data); err != nil {
		t.Error(err)
	}

	if err := s.SetChatMsgAsRead(data[0].FromId, data[0].ToId); err != nil {
		t.Error(err)
	}

	// chat message status=1
	get, err := s.GetChatMessageById(data[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if get.Status == 0 {
		t.Fatal("chat message status != 1")
	}

	// chat message status=0
	get, err = s.GetChatMessageById(data[2].ID)
	if err != nil {
		t.Fatal(err)
	}
	if get.Status == 1 {
		t.Fatal("chat message status != 0")
	}
}


func TestGetChatUserList(t *testing.T) {
	s := beforeChat()
	defer s.Close()

	data := []model.Chat {
		{
			FromId: 1,
			ToId:   0,
			Content: "Hello",
		},
		{
			FromId: 0,
			ToId:   1,
			Content: "Hello+1",
		},
		{
			FromId: 2,
			ToId:   0,
			Content: "Hello+2",
		},
	}
	if err := createChatMessageList(s, data); err != nil {
		t.Fatal(err)
	}

	gets, err := s.GetChatUserList(0, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != 2 {
		t.Fatal("get chat user err, size get:", len(gets), "expect:", 2)
	}

	// extra test
	user0 := gets[0]
	if user0.FromId != data[2].FromId && user0.Content != data[2].Content {
		t.Fatal("chat user[0] err, get:", user0, " expect:", data[2])
	}

	user1 := gets[1]
	if user1.FromId != data[1].FromId && user1.Content != data[1].Content {
		t.Fatal("chat user[1] err, get:", user1, " expect:", data[1])
	}
}


// 消息有可能是A发给B，或者B发给A，from, to 互相转换的
// 应该把两种情况都考虑到
func TestGetMsgListFromTo(t *testing.T) {
	s := beforeChat()
	defer s.Close()

	data := []model.Chat {
		{
			FromId: 1,
			ToId: 2,
			Content: "Hello",
		},
		{
			FromId: 2,
			ToId: 1,
			Content: "World",
		},
	}
	if err := createChatMessageList(s, data); err != nil {
		t.Error(err)
	}

	gets, err := s.GetChatMsgList(1, 2, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(gets) != 2 {
		t.Fatal("get msg list size err, get:", len(gets), "expect:", len(data))
	}
	m1 := gets[0]
	if m1.Content != data[1].Content {
		t.Fatal("get msg list 2-1 err")
	}
	m2 := gets[1]
	if m2.Content != data[0].Content {
		t.Fatal("get msg list 1-2 err")
	}
}

func TestGetChatId(t *testing.T) {
	var (
		from, to int64 = 1, 2
	)
	if getChatId(from, to) != getChatId(to, from) {
		t.Error("err!")
	}
	cid := getChatId(from, to)
	if cid != 0x0100000002 {
		t.Error("err!")
	}
}

func createChatMessageList(s *datasource, data []model.Chat) error {
	for i := range data {
		if err := s.CreateChatMessage(&data[i]); err != nil {
			return err
		}
	}
	return nil
}

func beforeChat() *datasource {
	s := newTest()
	s.Exec("DELETE FROM chats;")
	return s
}
