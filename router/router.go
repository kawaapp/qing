package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/server"
)

func Load(mwx ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(mwx...)

	// user register/login/logout
	e.POST("/api/login", server.HandleLogin)
	e.POST("/api/logout", server.HandleLogout)
	e.POST("/api/auth", server.HandleAuthMp)
	e.POST("/api/auth/mp", server.HandleAuthMp)
	e.POST("/api/auth/h5", server.HandleAuthH5)
	e.POST("/api/register", server.Register)


	// api - user
	get := e.Group("/api")
	{
		// user
		get.GET("/users/:id", server.GetUser)
		get.GET("/users/self", server.Self)
		get.GET("/users/:id/discussions", server.GetDiscussionByUser)
		get.GET("/users/:id/posts", server.GetPostByUser)
		get.GET("/users/:id/likes", server.GetLikeByUser)

		// discussion
		get.GET("/discussions", server.GetDiscussionList)
		get.GET("/discussions/:id", server.GetDiscussion)

		// posts
		get.GET("/discussions/:pid/posts", server.GetPostList)
		get.GET("/discussions/posts/:id", server.GetPost)

		// likes
		get.GET("/posts/:id/likes", server.GetLikeList)
		get.GET("/posts/:id/likes/count", server.GetLikeCount)

		// tags
		get.GET("/tags/:tag/discussions", server.GetDiscussionsByTag)
		get.GET("/tags", server.GetTagList)

		// notification - comment/mention/like
		get.GET("/notifications/count", server.GetNotificationCount)
		get.GET("/notifications", server.GetNotificationList)

		// chat
		get.GET("/chat/messages", server.GetChatListByUser)
		get.GET("/chat/users", server.GetChatUserList)
	}

	write := e.Group("/api")
	{
		write.Use(session.AttachUser())

		// user
		write.PUT("/users", server.UpdateUser)

		// discussion
		write.POST("/discussions", server.CreateDiscussion)
		write.DELETE("/discussions/:id", server.DeleteDiscussion)

		// posts
		write.POST("/discussions/posts", server.CreatePost)
		write.DELETE("/discussions/posts/:id", server.DeletePost)

		// likes
		write.POST("/posts/likes", server.CreateLike)
		write.DELETE("/posts/:id/likes", server.DeleteLike)

		// notification - comment/mention/like
		write.PUT("/notifications/:id/read", server.SetNotificationRead)
		write.PUT("/notifications/read", server.SetNotificationReadType)

		// chat
		write.POST("/chat/messages", server.CreateChatMessage)
		write.PUT("/chat/messages/read", server.SetChatMessageRead)

		// report
		write.POST("/reports", server.CreateReport)
	}

	// api - admin
	admin := e.Group("/api/b")
	{
		// stats
		admin.GET("/stats/overview", server.GetStatsOverView)
		admin.GET("/stats", server.GetStats)

		// users
		admin.GET("/users/search", server.SearchUser)
		admin.DELETE("/users/:id", server.DeleteUser)
		admin.PUT("/users/:id/st", server.SetUserStatus)

		// discussions
		admin.GET("/discussions/search", server.SearchDiscussions)
		admin.DELETE("/discussions/:id", server.DeleteDiscussion)
		admin.PUT("/discussions/:id/st", server.SetDiscussionStatus)
		admin.POST("/discussions", server.AdminCreateDiscussion)

		// posts
		admin.GET("/posts/search", server.SearchPosts)
		admin.DELETE("/posts/:id", server.DeletePost)
		admin.POST("/posts", server.CreatePost)

		// topic (it's hash tag)
		admin.GET("/tags", server.GetTagList)
		admin.DELETE("/tags/:id", server.DeleteTag)
		admin.PUT("/tags/:id", server.UpdateTag)
		admin.POST("/tags", server.CreateTag)
		admin.POST("/tags/posts", server.LinkTagPost)

		// report
		admin.GET("/reports/search", server.SearchReport)
		admin.GET("/reports/:id", server.GetReport)
		admin.PUT("/reports/:id/st", server.SetReportStatus)

		// spam check
		admin.GET("/spam/words", server.GetSpamWords)
		admin.POST("/spam/words", server.CreateSpamWords)
		admin.POST("/spam/check", server.SpamCheck)
	}
	return e
}
