package server

const (
	eDiscussionCreated = "discussion.created"
	eDiscussionDeleted = "discussion.deleted"

	ePostCreated = "post.created"
	ePostDeleted = "post.deleted"

	eLikeCreated  = "favor.created"
	eLikeDeleted  = "favor.deleted"
	eLikeUpdated  = "favor.updated"

	ePostShared    = "post.share"
	ePostValued    = "post.value"
	ePostTopped    = "post.top"

	// 用户登录
	eUserLogin  = "user.login"
	eUserFollow = "user.follow"
	eUserUnfollow = "user.unfollow"
)