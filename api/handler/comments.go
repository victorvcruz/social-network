package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"social_network_project/controllers"
	"social_network_project/controllers/errors"
	"social_network_project/controllers/validate"
	"social_network_project/entities"
	"time"
)

type CommentsAPI struct {
	CommentController controllers.CommentsController
	Validate          *validator.Validate
}

func RegisterCommentsHandlers(handler *gin.Engine, commentsController controllers.CommentsController) {
	ac := &CommentsAPI{
		CommentController: commentsController,
		Validate:          validator.New(),
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
	postID := c.Param("post")
	commentID := c.DefaultQuery("comment_id", "")

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	comment := CreateCommentStruct(mapBody, accountID, &postID, &commentID)

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(comment)
	if err != nil {
		mapper["errors"] = validate.RequestCommentValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	err = a.CommentController.InsertComment(comment)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundPostIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}
	c.JSON(http.StatusOK, comment.ToResponse())
	return
}

func (a *CommentsAPI) GetComment(c *gin.Context) {

	accountID, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	postID := c.DefaultQuery("post_id", "")
	commentID := c.DefaultQuery("comment_id", "")

	comments, err := a.CommentController.FindCommentsByAccountID(accountID, &postID, &commentID)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundPostIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, comments)
	return
}

func (a *CommentsAPI) UpdateComment(c *gin.Context) {
	accountID, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}
	postID := ""
	commentID := ""

	comment := CreateCommentStruct(mapBody, accountID, &postID, &commentID)
	comment.ID = stringNullable(mapBody["id"])

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(comment)
	if err != nil {
		mapper["errors"] = validate.RequestCommentValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	commentUpdated, err := a.CommentController.UpdateCommentDataByID(comment)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, commentUpdated)
	return
}

func (a *CommentsAPI) DeleteComment(c *gin.Context) {
	accountID, err := decodeTokenAndReturnID(c.Request.Header.Get("BearerToken"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}
	postID := ""
	commentID := ""
	comment := CreateCommentStruct(mapBody, accountID, &postID, &commentID)
	comment.ID = stringNullable(mapBody["id"])
	comment.Content = "--"

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(comment)
	if err != nil {
		mapper["errors"] = validate.RequestCommentValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	commentToRemoved, err := a.CommentController.RemoveCommentByID(comment, accountID)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, commentToRemoved)
	return

}

func CreateCommentStruct(mapBody map[string]interface{}, accountID, postID, commentID *string) *entities.Comment {

	return &entities.Comment{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		PostID:    *postID,
		CommentID: *commentID,
		Content:   stringNullable(mapBody["content"]),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

}
