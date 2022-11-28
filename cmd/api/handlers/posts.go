package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"social_network_project/internal/platform/cache"
	post2 "social_network_project/internal/post"
	"social_network_project/internal/post/service"
	"social_network_project/internal/utils"
	"social_network_project/internal/utils/errors"
	"social_network_project/internal/utils/validate"
	"strconv"
	"time"
)

type PostHandlerClient interface {
	CreatePost(c *gin.Context)
	GetPost(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
	SearchPostByAccountFollowing(c *gin.Context)
}

type PostsAPI struct {
	Controller  service.PostsServiceClient
	RedisClient cache.RedisServiceClient
	Validate    *validator.Validate
}

func RegisterPostsHandlers(postsController service.PostsServiceClient, _redis cache.RedisServiceClient) PostHandlerClient {
	return &PostsAPI{
		Controller:  postsController,
		RedisClient: _redis,
		Validate:    validator.New(),
	}
}

func (a *PostsAPI) CreatePost(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request post2.PostRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}
	post := a.fillFields(request, &accountID)

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(post)
	if err != nil {
		mapper["errors"] = validate.RequestPostValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	err = a.Controller.InsertPost(post)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
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

	idToGet := c.DefaultQuery("account_id", accountID)
	page := c.DefaultQuery("page", "1")
	if _, err = strconv.ParseInt(page, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Page is not a number",
		})
		return
	}

	postsOfAccount, err := a.Controller.FindPostsByAccountID(&accountID, &idToGet, &page)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	a.RedisClient.InsertCache(c.Request, postsOfAccount)
	c.JSON(http.StatusOK, postsOfAccount)
	return
}

func (a *PostsAPI) UpdatePost(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request post2.PostRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}
	post := a.fillFields(request, &accountID)
	post.ID = utils.StringNullable(request.Id)

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(post)
	if err != nil {
		mapper["errors"] = validate.RequestPostValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	postUpdated, err := a.Controller.UpdatePostDataByID(post)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundPostIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
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
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request post2.PostRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}
	post := a.fillFields(request, &accountID)
	post.ID = utils.StringNullable(request.Id)
	post.Content = "--"

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(post)
	if err != nil {
		mapper["errors"] = validate.RequestPostValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	postToRemoved, err := a.Controller.RemovePostByID(post)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundPostIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, postToRemoved)
	return
}

func (a *PostsAPI) SearchPostByAccountFollowing(c *gin.Context) {

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

	page := c.DefaultQuery("page", "1")
	if _, err = strconv.ParseInt(page, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Page is not a number",
		})
		return
	}
	postsOfAccount, err := a.Controller.FindPostByAccountFollowingByAccountID(&accountID, &page)
	if err != nil {
		switch e := err.(type) {
		case *errors.NotFoundAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	a.RedisClient.InsertCache(c.Request, postsOfAccount)
	c.JSON(http.StatusOK, postsOfAccount)
	return
}

func (a *PostsAPI) fillFields(req post2.PostRequest, accountID *string) *post2.Post {

	return &post2.Post{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		Content:   utils.StringNullable(req.Content),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}
}
