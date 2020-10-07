package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/events"

	"net/http"
	"strconv"
)

// 返回消息汇总数据
func GetNotificationCount(c echo.Context) error {
	usr := session.User(c)
	mc, err := store.GetNotificationCount(c, usr.ID)
	if err != nil {
		return err
	}
	return c.JSON(200, mc)
}

func GetNotificationList(c echo.Context) error {
	var (
		q   = c.FormValue("q")
		usr = session.User(c)
		err error
	)

	var messages []*model.Notification
	var mt = model.NotUnknown
	if q == "favor" {
		mt = model.NotLike
	} else if q == "comment" {
		mt = model.NotComment
	}

	page, size := getPageSize(c)
	messages, err = store.GetNotificationListType(c, usr.ID, mt, page-1, size)
	if err != nil {
		return err
	}

	p := makePayload(0, messages)

	// attach user
	if includes(c, "user") {
		attachUserToNotification(c, messages, p.Entities)
	}

	return c.JSON(200, messages)
}

func SetNotificationRead(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "message id not found")
	}
	uid := session.User(c).ID
	err = store.SetNotificationReadId(c, uid, int64(id))
	if err != nil {
		return err
	}
	return c.NoContent(200)
}

func SetNotificationReadType(c echo.Context) error {
	var (
		typo  = c.QueryParam("type")
		err error
	)

	if typo == "like" {
		err = store.SetNotificationReadType(c, session.User(c).ID, model.NotLike)
	} else if typo == "comment" {
		err = store.SetNotificationReadType(c, session.User(c).ID, model.NotComment)
	}
	if err != nil {
		return err
	}
	return c.NoContent(200)
}


func attachUserToNotification(c echo.Context, messages []*model.Notification, entities map[string]interface{}) {
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

// events
func msgOnCommentCreated(c echo.Context, v interface{}) error {
	comment, ok := v.(*model.Post)
	if !ok {
		return typeError("Comment")
	}

	n := &model.Notification {
		EntityId: comment.ID,
		EntityType: int(model.NotComment),
		FromId: comment.AuthorID,
		ToId: comment.ReplyID,
	}
	return store.CreateNotification(c, n)
}

//
func msgOnCommentDeleted(c echo.Context, v interface{}) error {
	return nil
}

//
func msgOnLikeCreated(c echo.Context, v interface{}) error {
	f, ok := v.(*model.Like)
	if !ok {
		return typeError("Like")
	}

	toId, err := getLikedUser(c, f)
	if err != nil {
		return err
	}

	n := &model.Notification {
		EntityId: f.ID,
		EntityType: int(model.NotLike),
		FromId: f.UserID,
		ToId: toId,
	}
	return store.CreateNotification(c, n)
}

//
func msgOnLikeUpdated(c echo.Context, v interface{}) error {
	return nil
}

func getLikedUser(c echo.Context, l *model.Like) (int64, error) {
	if l.TargetTy == model.LikeDiscussion {
		d, err := store.GetDiscussion(c, l.TargetID)
		if err != nil {
			return 0, err
		}
		return d.AuthorID, nil
	} else {
		p, err := store.GetPost(c, l.TargetID)
		if err != nil {
			return 0, err
		}
		return p.AuthorID, nil
	}
}

func init() {
	events.Subscribe(ePostCreated, msgOnCommentCreated)
	events.Subscribe(ePostDeleted, msgOnCommentDeleted)
	events.Subscribe(eLikeCreated, msgOnLikeCreated)
	events.Subscribe(eLikeUpdated, msgOnLikeUpdated)
}