package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"social_network_project/internal/comment"
	"social_network_project/internal/comment/service"
	"social_network_project/internal/platform/cache"
	"social_network_project/internal/utils"
	"social_network_project/internal/utils/errors"
	"social_network_project/internal/utils/validate"
	"strconv"
	"time"
)

type CommentsHandlerClient interface {
	CreateComment(c *gin.Context)
	GetComment(c *gin.Context)
	UpdateComment(c *gin.Context)
	DeleteComment(c *gin.Context)
}

type CommentsHandler struct {
	Controller  service.CommentsServiceClient
	RedisClient cache.RedisServiceClient
	Validate    *validator.Validate
}

func RegisterCommentsHandlers(commentsController service.CommentsServiceClient, _redis cache.RedisServiceClient) CommentsHandlerClient {
	return &CommentsHandler{
		Controller:  commentsController,
		RedisClient: _redis,
		Validate:    validator.New(),
	}
}

func (a *CommentsHandler) CreateComment(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}
	postID := c.Param("post")
	commentID := c.DefaultQuery("comment_id", "")

	var request comment.CommentRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	comment := a.fillFields(request, accountID, postID, utils.NewNullString(commentID))

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(comment)
	if err != nil {
		mapper["errors"] = validate.RequestCommentValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	err = a.Controller.InsertComment(comment)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundPostIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}
	c.JSON(http.StatusOK, comment.ToResponse())
	return
}

func (a *CommentsHandler) GetComment(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	resCache, err := a.RedisClient.FindInCache(c.Request)
	switch e := err.(type) {
	case *errors.CacheNotFoundError:
		log.Println(e.Error())
	default:
		c.JSON(http.StatusOK, resCache)
		log.Println("Cache")

		return
	}

	idToGet := c.DefaultQuery("account_id", "")
	postID := c.DefaultQuery("post_id", "")
	commentID := c.DefaultQuery("comment_id", "")
	page := c.DefaultQuery("page", "1")
	if _, err = strconv.ParseInt(page, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Page is not a number",
		})
		return
	}

	comments, err := a.Controller.FindCommentsByAccountID(&accountID, &idToGet, &postID, &commentID, &page)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundPostIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	a.RedisClient.InsertCache(c.Request, comments)
	c.JSON(http.StatusOK, comments)
	return
}

func (a *CommentsHandler) UpdateComment(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request comment.CommentRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	comment := a.fillFields(request, accountID, "", utils.NewNullString(""))
	comment.ID = utils.StringNullable(request.Id)

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(comment)
	if err != nil {
		mapper["errors"] = validate.RequestCommentValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	commentUpdated, err := a.Controller.UpdateCommentDataByID(comment)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, commentUpdated)
	return
}

func (a *CommentsHandler) DeleteComment(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request comment.CommentRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	comment := a.fillFields(request, accountID, "", utils.NewNullString(""))
	comment.ID = utils.StringNullable(request.Id)
	comment.Content = "--"

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(comment)
	if err != nil {
		mapper["errors"] = validate.RequestCommentValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	commentToRemoved, err := a.Controller.RemoveCommentByID(comment)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundCommentIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, commentToRemoved)
	return

}

func (a *CommentsHandler) fillFields(req comment.CommentRequest, accountID, postID string, commentID sql.NullString) *comment.Comment {

	return &comment.Comment{
		ID:        uuid.New().String(),
		AccountID: accountID,
		PostID:    postID,
		CommentID: commentID,
		Content:   utils.StringNullable(req.Content),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

}
