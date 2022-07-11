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

type PostsAPI struct {
	PostController controllers.PostsController
	Validate       *validator.Validate
}

func RegisterPostsHandlers(handler *gin.Engine, postsController controllers.PostsController) {
	ac := &PostsAPI{
		PostController: postsController,
		Validate:       validator.New(),
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

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	post := CreatePostStruct(mapBody, accountID)

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(post)
	if err != nil {
		mapper["errors"] = validate.RequestPostValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	err = a.PostController.InsertPost(post)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
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

	postsOfAccount, err := a.PostController.FindPostsByAccountID(accountID)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
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

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	post := CreatePostStruct(mapBody, accountID)
	post.ID = stringNullable(mapBody["id"])

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(post)
	if err != nil {
		mapper["errors"] = validate.RequestPostValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	postUpdated, err := a.PostController.UpdatePostDataByID(post)
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
		default:
			log.Fatal(err)
		}
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

	mapBody, err := readBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	post := CreatePostStruct(mapBody, accountID)
	post.ID = stringNullable(mapBody["id"])
	post.Content = "--"

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(post)
	if err != nil {
		mapper["errors"] = validate.RequestPostValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	postToRemoved, err := a.PostController.RemovePostByID(post, accountID)
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
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, postToRemoved)
	return
}

func CreatePostStruct(mapBody map[string]interface{}, accountID *string) *entities.Post {

	return &entities.Post{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		Content:   stringNullable(mapBody["content"]),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}
}

func stringNullable(str interface{}) string {
	if str == nil {
		return ""
	}
	return str.(string)
}
