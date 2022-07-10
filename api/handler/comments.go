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

type CommentsAPI struct {
	CommentController controllers.CommentsController
	PostController    controllers.PostsController
	AccountController controllers.AccountsController
}

func RegisterCommentsHandlers(handler *gin.Engine, commentsController controllers.CommentsController, postController controllers.PostsController, accountsController controllers.AccountsController) {
	ac := &CommentsAPI{
		CommentController: commentsController,
		PostController:    postController,
		AccountController: accountsController,
	}

	handler.POST("/comments/:post", ac.CreateComment)
	handler.GET("/accounts/comments", ac.GetComment)
	handler.PUT("/comments", ac.UpdateComment)
	handler.DELETE("/comments", ac.DeleteComment)
}

func (a *CommentsAPI) CreateComment(c *gin.Context) {
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

	comment := CreateCommentStruct(mapBody, c.Param("post"), c.DefaultQuery("comment_id", ""), accountID)

	err = a.CommentController.InsertComment(comment)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, comment.ToResponse())
	return
}

func (a *CommentsAPI) GetComment(c *gin.Context) {

}

func (a *CommentsAPI) UpdateComment(c *gin.Context) {

}

func (a *CommentsAPI) DeleteComment(c *gin.Context) {

}

func CreateCommentStruct(mapBody map[string]interface{}, postID string, commentID string, accountID *string) *entities.Comment {

	return &entities.Comment{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		PostID:    postID,
		CommentID: commentID,
		Content:   mapBody["content"].(string),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

}
