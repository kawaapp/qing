package store

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"time"
)

type UserStore interface {
	GetUserList(params model.QueryParams, page, size int) ([]*model.User, error)
	GetUserCount(params model.QueryParams) (int, error)

	GetUser(id int64) (*model.User, error)
	GetUserIdList(ids []int64) ([]*model.User, error)
	GetUserByLogin(string) (*model.User, error)
	GetUserByPhone(string) (*model.User, error)
	CreateUser(usr *model.User) error
	UpdateUser(usr *model.User) error
	DeleteUser(id int64) error

	// bind
	GetUserByBind(string) (*model.User, error)
	GetUserByUnion(string) (*model.User, error)
	CreateBind(bind *model.UserBind) error
	DeleteBind(bindId string) error
	DeleteBindByUser(uid int64) error
	DeleteBindUser(uid int64) error
	CreateBindUser(user *model.User, bind *model.UserBind) error

	// user-exp-sign
	UpdateUserExp(uid int64, exp int) error
	UpdateUserSign(uid int64, count int) error

	// follow
	GetFollowerList(uid int64, page, size int) ([]*model.User, error)
	GetFollowingList(uid int64, page, size int) ([]*model.User, error)
	GetFollow(uid, fid int64) (*model.Follow, error)

	CreateFollow(f *model.Follow) error
	DeleteFollow(uid, fid int64) error
}

type DiscussionStore interface {
	GetDiscussionList(params model.QueryParams, page, size int) ([]*model.Discussion, error)
	GetDiscussionCount(params model.QueryParams) (int, error)

	GetDiscussionListByIds(ids []int64) ([]*model.Discussion, error)
	GetDiscussion(id int64) (*model.Discussion, error)
	CreateDiscussion(p *model.Discussion) error
	UpdateDiscussion(p *model.Discussion) error
	DeleteDiscussion(id int64) error
}

type FavoriteStore interface {
	GetFavoriteList(q model.QueryParams, page, size int) ([]*model.Favorite, error)
	GetFavoriteCount(q model.QueryParams) (int, error)

	CreateFavorite(f *model.Favorite) error
	GetFavoriteUser(uid, pid int64) (*model.Favorite, error)
	GetFavoriteId(id int64) (*model.Favorite, error)
	DeleteFavorite(id int64) error
}

type PostStore interface {
	GetPostList(params model.QueryParams, page, size int) ([]*model.Post, error)
	GetPostCount(params model.QueryParams) (int, error)

	GetPostListUser(uid int64, page, size int) ([]*model.Post, error)
	GetPostListByIds(ids []int64) ([]*model.Post, error)
	GetPost(id int64) (*model.Post, error)

	CreatePost(p *model.Post) error
	UpdatePost(p *model.Post) error
	DeletePost(id int64) error
}

type LikeStore interface {
	GetLikeList(q model.QueryParams, page, size int) ([]*model.Like, error)
	GetLikeCount(q model.QueryParams) (int, error)

	GetLike(ty string, id, uid int64) (*model.Like, error)
	GetLikeId(id int64) (*model.Like, error)

	CreateLike(t string, tid, uid int64) (error, bool)
	DeleteLike(t string, tid, uid int64) error

	// 从数组中过滤用户赞过的Post
	GetUserLikedPostList(uid int64, pids []int64) ([]int64, error)
}

type CategoryStore interface {
	GetCategoryList() ([]*model.Category, error)
	GetCategory(int64) (*model.Category, error)
	CreateCategory(*model.Category) error
	UpdateCategory(*model.Category) error
	DeleteCategory(int64) error
}

type TagStore interface {
	GetDiscussionsByTag(tag string, page, size int) ([]*model.Discussion, error)
	GetTagList() ([]*model.Tag, error)
	GetTagId(id int64) (*model.Tag, error)
	LinkTagDiscussion(pid int64, tag []string) error
	CreateTag(tag, summary string) (*model.Tag, error)
	DeleteTag(id int64) error
	UpdateTag(t *model.Tag) error
}

type MessageStore interface {
	// notification
	GetNotificationCount(to int64) (*model.MessageCount, error)
	GetNotificationById(id int64) (*model.Notification, error)
	GetNotificationListType(to int64, mt model.NotType, page, size int) ([]*model.Notification, error)
	CreateNotification(n *model.Notification) error
	SetNotificationReadId(uid, id int64) error
	SetNotificationReadType(to int64, mt model.NotType) error

	// chat
	GetChatMsgList(from, to int64, page, size int) ([]*model.Chat, error)
	GetChatMessageById(id int64) (*model.Chat, error)
	CreateChatMessage(m *model.Chat) error
	SetChatMsgAsRead(from, to int64) error
	GetChatUserList(to int64, page, size int) ([]*model.Chat, error)
}

type MetaStore interface {
	GetMetaData() (map[string]string, error)
	GetMetaValue(key string) (string, error)
	PutMetaData(map[string]interface{}) error
}

type MediaStore interface {
	CreateMedia(m *model.Media) error
	GetMediaListByPostIds(pids []int64) ([]*model.Media, error)
	GetMediaByPostId(pid int64) (*model.Media, error)
	DeleteMediaByPostId(pid int64) error
}

type ReportStore interface {
	GetReportList(q model.QueryParams, page, size int) ([]*model.Report, error)
	GetReportCount(q model.QueryParams) (int, error)

	CreateReport(rpt *model.Report) error
	GetReport(id int64) (*model.Report, error)
	UpdateReport(*model.Report) error
}

type AnalyseStore interface {
	// analytics
	GetTotalUser() (int, error)
	GetNewUser(day time.Time) (int, error)
	GetNewUserDaily(from, to time.Time) ([]*model.DailyCount, error)

	GetUserActive(day time.Time) (int, error)
	GetUserActiveDaily(from, to time.Time) ([]*model.DailyCount, error)

	// analytics
	GetTotalDiscussion() (int, error)
	GetNewDiscussion(day time.Time) (int, error)
	GetNewDiscussionDaily(from, to time.Time) ([]*model.DailyCount, error)
}

type Store interface {
	// user
	UserStore

	// discussion
	DiscussionStore

	// posts
	PostStore

	// like
	LikeStore

	// category
	CategoryStore

	// tag
	TagStore

	// message
	MessageStore

	// meta-data
	MetaStore

	// images
	MediaStore

	// report
	ReportStore

	// analytics
	AnalyseStore

	// favorite
	FavoriteStore
}

// user
func GetUser(c echo.Context, id int64) (*model.User, error) {
	return FromContext(c).GetUser(id)
}

func GetUserLogin(c echo.Context, login string) (*model.User, error) {
	return FromContext(c).GetUserByLogin(login)
}

func GetUserPhone(c echo.Context, phone string) (*model.User, error)  {
	return FromContext(c).GetUserByPhone(phone)
}

func GetUserBind(c echo.Context, bind string) (*model.User, error) {
	return FromContext(c).GetUserByBind(bind)
}

func GetUserUnion(c echo.Context, bind string) (*model.User, error) {
	return FromContext(c).GetUserByUnion(bind)
}

func CreateUser(c echo.Context, usr *model.User) error {
	return FromContext(c).CreateUser(usr)
}

func CreateBind(c echo.Context, bind *model.UserBind) error {
	return FromContext(c).CreateBind(bind)
}

func DeleteBind(c echo.Context, bindId string) error {
	return FromContext(c).DeleteBind(bindId)
}

func UpdateUser(c echo.Context, usr *model.User) error {
	return FromContext(c).UpdateUser(usr)
}

func DeleteUser(c echo.Context, id int64) error {
	return FromContext(c).DeleteUser(id)
}

// discussion
func GetDiscussion(c echo.Context, id int64) (*model.Discussion, error) {
	return FromContext(c).GetDiscussion(id)
}

func CreateDiscussion(c echo.Context, p *model.Discussion) error {
	return FromContext(c).CreateDiscussion(p)
}

func UpdateDiscussion(c echo.Context, p *model.Discussion) error {
	return FromContext(c).UpdateDiscussion(p)
}

func DeleteDiscussion(c echo.Context, id int64) error {
	return FromContext(c).DeleteDiscussion(id)
}

// posts
func GetPost(c echo.Context, id int64) (*model.Post, error) {
	return FromContext(c).GetPost(id)
}

func CreatePost(c echo.Context, cmt *model.Post) error {
	return FromContext(c).CreatePost(cmt)
}

func UpdatePost(c echo.Context, cmt *model.Post) error {
	return FromContext(c).UpdatePost(cmt)
}

func DeletePost(c echo.Context, id int64) error {
	return FromContext(c).DeletePost(id)
}

// like
func GetLikeList(c echo.Context, q model.QueryParams, page, size int) ([]*model.Like, error) {
	return FromContext(c).GetLikeList(q, page, size)
}

func GetLikeCount(c echo.Context, q model.QueryParams) (int, error) {
	return FromContext(c).GetLikeCount(q)
}

func GetLike(c echo.Context, t string, tid, uid int64) (*model.Like, error) {
	return FromContext(c).GetLike(t, tid, uid)
}

func CreateLike(c echo.Context, t string, tid, uid int64) (error,bool) {
	return FromContext(c).CreateLike(t, tid, uid)
}

func DeleteLike(c echo.Context, t string, tid, uid int64) error {
	return FromContext(c).DeleteLike(t, tid, uid)
}

// tags
func GetTagList(c echo.Context) ([]*model.Tag, error) {
	return FromContext(c).GetTagList()
}

func GetDiscussionByTag(c echo.Context, tag string, page, size int) ([]*model.Discussion, error) {
	return FromContext(c).GetDiscussionsByTag(tag, page, size)
}

func LinkTagPost(c echo.Context, pid int64, tags []string) error {
	return FromContext(c).LinkTagDiscussion(pid, tags)
}

// message
func GetNotificationCount(c echo.Context, to int64) (*model.MessageCount, error) {
	return FromContext(c).GetNotificationCount(to)
}

func GetNotificationListType(c echo.Context, to int64, mt model.NotType, page, size int) ([]*model.Notification, error) {
	return FromContext(c).GetNotificationListType(to, mt, page, size)
}

func SetNotificationReadId(c echo.Context, uid, id int64) error {
	return FromContext(c).SetNotificationReadId(uid, id)
}

func GetNotificationId(c echo.Context, id int64) (*model.Notification, error) {
	return FromContext(c).GetNotificationById(id)
}

func SetNotificationReadType(c echo.Context, to int64, mt model.NotType) error {
	return FromContext(c).SetNotificationReadType(to, mt)
}

func SetNotificationReadFromTo(c echo.Context, from, to int64) error {
	return FromContext(c).SetChatMsgAsRead(from, to)
}

func CreateNotification(c echo.Context, n *model.Notification) error {
	return FromContext(c).CreateNotification(n)
}

// meta
func GetMetaValue(c echo.Context, key string) (string, error) {
	return FromContext(c).GetMetaValue(key)
}

// counter

func GetDailyActiveUser(c echo.Context, from, to time.Time) ([]*model.DailyCount, error) {
	return FromContext(c).GetUserActiveDaily(from, to)
}

func GetDailyNewUser(c echo.Context, from, to time.Time) ([]*model.DailyCount, error) {
	return FromContext(c).GetNewUserDaily(from, to)
}

func GetDailyNewDiscussion(c echo.Context, from, to time.Time) ([]*model.DailyCount, error)  {
	return FromContext(c).GetNewDiscussionDaily(from, to)
}

func GetDailyNewComment(c echo.Context, from, to time.Time) ([]*model.DailyCount, error) {
	return nil, nil
}

func GetDailyNewFavor(c echo.Context, from, to time.Time) ([]*model.DailyCount, error)  {
	return nil, nil
}
