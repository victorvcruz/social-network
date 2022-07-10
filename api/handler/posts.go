package handler

import (
	"github.com/gin-gonic/gin"
	"social_network_project/controllers"
)

type PostsAPI struct {
	Controller controllers.PostsController
}

func RegisterPostsHandlers(handler *gin.Engine, postsController controllers.PostsController) {
	ac := &PostsAPI{
		Controller: postsController,
	}

	handler.POST("/posts", ac.CreatePost)
	handler.PUT("/posts", ac.UpdatePost)
	handler.DELETE("/posts", ac.DeletePost)
}

func (a *PostsAPI) CreatePost(c *gin.Context) {

}

func (a *PostsAPI) UpdatePost(c *gin.Context) {

}

func (a *PostsAPI) DeletePost(c *gin.Context) {

}
