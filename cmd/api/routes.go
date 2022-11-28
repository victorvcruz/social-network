package api

import (
	"github.com/gin-gonic/gin"
	"social_network_project/cmd/api/handlers"
)

func Init(
	auth handlers.AuthHandlerClient,
	accounts handlers.AccountsHandlerClient,
	posts handlers.PostHandlerClient,
	comments handlers.CommentsHandlerClient,
	intercations handlers.IntercationsHandlerClient,
	) *gin.Engine {
	app := gin.Default()

	app.POST("/auth", auth.CreateToken)

	app.POST("/accounts", accounts.CreateAccount)
	app.GET("/accounts", accounts.GetAccount)
	app.PUT("/accounts", accounts.UpdateAccount)
	app.DELETE("/accounts", accounts.DeleteAccount)
	app.POST("/accounts/follows", accounts.FollowAccount)
	app.GET("/accounts/following", accounts.SearchFollowing)
	app.GET("/accounts/follower", accounts.SearchFollowers)
	app.DELETE("/accounts/follows", accounts.UnfollowAccount)

	app.POST("/comments/:post", comments.CreateComment)
	app.GET("/accounts/comments", comments.GetComment)
	app.PUT("/comments", comments.UpdateComment)
	app.DELETE("/comments", comments.DeleteComment)

	app.POST("/interaction", intercations.CreateInteraction)
	app.PUT("/interaction", intercations.UpdateInteraction)
	app.DELETE("/interaction", intercations.DeleteInteraction)

	app.POST("/posts", posts.CreatePost)
	app.GET("/accounts/posts", posts.GetPost)
	app.PUT("/posts", posts.UpdatePost)
	app.DELETE("/posts", posts.DeletePost)
	app.GET("/accounts/follows/posts", posts.SearchPostByAccountFollowing)

	return app
}
