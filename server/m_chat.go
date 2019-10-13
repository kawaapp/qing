package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/model"

	"strconv"
	"strings"
)

type chatMessage struct {
	From int64 `json:"from_id"`
	To int64 `json:"to_id"`
	Content string `json:"content"`
	Type int `json:"type"`
}

func GetChatUserList(c echo.Context) error {
	var (
		uid = session.User(c).ID
		page, limit = getPageSize(c)
	)
	data, err := store.FromContext(c).GetChatUserList(uid, page, limit)
	if err != nil {
		return err
	}

	p := makePayload(0, data)

	// attach user
	if strings.Contains(c.QueryParam("includes"), "user") {
		attachUserToChat(c, data, p.Entities)
	}
	return c.JSON(200, data)
}

// 当调用这个方法的时候, 可以直接把收到的消息标记为已读
func GetChatListByUser(c echo.Context) error {
	from, err := strconv.Atoi(c.QueryParam("from"))
	if err != nil {
		return c.String(400, "from id not found")
	}
	var (
		uid = session.User(c).ID
		page, limit = getPageSize(c)
	)
	data, err := store.FromContext(c).GetChatMsgList(int64(from), uid, page, limit)
	if err != nil {
		return err
	}

	p := makePayload(0, data)

	// attach user
	if strings.Contains(c.QueryParam("includes"), "user") {
		attachUserToChat(c, data, p.Entities)
	}

	// mark as read
	go func(store store.Store) {
		store.SetChatMsgAsRead(int64(from), uid)
	}(store.FromContext(c))
	return c.JSON(200, p)
}

func CreateChatMessage(c echo.Context) error {
	in := new(chatMessage)
	if err := c.Bind(&in); err != nil {
		return err
	}

	var uid int64
	if user := session.User(c); user.Login == "Admin" {
		uid = in.To
	} else {
		uid = user.ID
	}
	chat := model.Chat {
		FromId: uid,
		ToId: in.To,
		Content: in.Content,
	}
	err := store.FromContext(c).CreateChatMessage(&chat)
	if err != nil {
		return err
	}
	return c.JSON(200, &chat)
}

func SetChatMessageRead(c echo.Context) (err error) {

	from, err := strconv.Atoi(c.QueryParam("from"))
	if err != nil {
		return c.String(400, "from id not found")
	}
	err = store.FromContext (c).SetChatMsgAsRead(int64(from), session.User(c).ID)

	if err != nil {
		return err
	}
	return c.NoContent(200)
}

func attachUserToChat(c echo.Context, messages []*model.Chat, entities map[string]interface{}) {
	Ids := make([]int64, len(messages))
	kv := make(map[int64]*model.User)
	usr := session.User(c)

	for i, v := range messages {
		if usr.ID == v.FromId {
			Ids[i] = v.ToId
		} else {
			Ids[i] = v.FromId
		}
	}
	actors, _ := store.FromContext(c).GetUserIdList(Ids)
	if actors != nil {
		for _, v := range actors {
			kv[v.ID] = v
		}
	}
	entities["users"] = kv
}
