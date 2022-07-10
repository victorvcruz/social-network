package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"social_network_project/controllers"
	"social_network_project/entities"
	"time"
)

type PostsAPI struct {
	PostsController   controllers.PostsController
	AccountController controllers.AccountsController
}

func RegisterPostsHandlers(handler *gin.Engine, postsController controllers.PostsController, accountsController controllers.AccountsController) {
	ac := &PostsAPI{
		PostsController:   postsController,
		AccountController: accountsController,
	}

	handler.POST("/posts", ac.CreatePost)
	handler.GET("/accounts/posts", ac.GetPost)
	handler.PUT("/posts", ac.UpdatePost)
	handler.DELETE("/posts", ac.DeletePost)
}

func (a *PostsAPI) CreatePost(c *gin.Context) {

	accountID, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.AccountController.ExistsAccountByID(accountID)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Account id does not exist",
		})
		return
	}

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	if mapBody["content"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Add content",
		})
		return
	}

	post := CreatePostStruct(mapBody, accountID)

	err = a.PostsController.InsertPost(post)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, post.ToResponse())
	return

}

func (a *PostsAPI) GetPost(c *gin.Context) {

	accountID, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.AccountController.ExistsAccountByID(accountID)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Account id does not exist",
		})
		return
	}

	postsOfAccount, err := a.PostsController.FindPostsByAccountID(accountID)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, postsOfAccount)
	return
}

func (a *PostsAPI) UpdatePost(c *gin.Context) {
	accountID, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.AccountController.ExistsAccountByID(accountID)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Account id does not exist",
		})
		return
	}

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	if mapBody["content"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Add content",
		})
		return
	}
	if mapBody["id"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Add id",
		})
		return
	}

	postID := mapBody["id"].(string)

	existID, err = a.PostsController.ExistsPostByID(&postID)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Post id does not exist",
		})
		return
	}

	a.PostsController.ChangePostDataByID(&postID, mapBody["content"].(string))

	postUpdated, err := a.PostsController.FindPostByID(&postID)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, postUpdated)
	return
}

func (a *PostsAPI) DeletePost(c *gin.Context) {
	accountID, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	existID, err := a.AccountController.ExistsAccountByID(accountID)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Account id does not exist",
		})
		return
	}

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	if mapBody["id"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Add id",
		})
		return
	}

	postID := mapBody["id"].(string)

	existID, err = a.PostsController.ExistsPostByID(&postID)
	if err != nil {
		log.Fatal(err)
	}
	if !*existID {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Post id does not exist",
		})
		return
	}

	postToRemoved, err := a.PostsController.FindPostByID(&postID)
	if err != nil {
		log.Fatal(err)
	}

	a.PostsController.RemovePostByID(&postID)

	c.JSON(http.StatusOK, postToRemoved)
	return
}

func CreatePostStruct(mapBody map[string]interface{}, accountID *string) *entities.Post {

	post := entities.Post{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		Content:   mapBody["content"].(string),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

	return &post

}
