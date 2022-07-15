package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"social_network_project/cache"
	"social_network_project/cache/redisDB"
	"social_network_project/controllers"
	"social_network_project/controllers/errors"
	"social_network_project/controllers/validate"
	"social_network_project/entities"
	"social_network_project/utils"
	"strconv"
	"time"
)

type PostsAPI struct {
	Controller  controllers.PostsController
	RedisClient *redisDB.RedisClient
	Validate    *validator.Validate
}

func RegisterPostsHandlers(handler *gin.Engine, postsController controllers.PostsController, redisDB *redisDB.RedisClient) {
	ac := &PostsAPI{
		Controller:  postsController,
		RedisClient: redisDB,
		Validate:    validator.New(),
	}

	handler.POST("/posts", ac.CreatePost)
	handler.GET("/accounts/posts", ac.GetPost)
	handler.PUT("/posts", ac.UpdatePost)
	handler.DELETE("/posts", ac.DeletePost)
	handler.GET("/accounts/follows/posts", ac.SearchPostByAccountFollowing)
}

func (a *PostsAPI) CreatePost(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
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

	err = a.Controller.InsertPost(post)
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

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	resCache, err := cache.FindInCache(c.Request, a.RedisClient)
	switch e := err.(type) {
	case *errors.CacheNotFoundError:
		log.Println(e.Error())
	default:
		c.JSON(http.StatusOK, resCache)
		log.Println("Cache")

		return
	}

	idToGet := c.DefaultQuery("account_id", "")
	page := c.DefaultQuery("page", "1")
	if _, err = strconv.ParseInt(page, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Page is not a number",
		})
		return
	}

	postsOfAccount, err := a.Controller.FindPostsByAccountID(accountID, &idToGet, &page)
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

	cache.InsertCache(c.Request, postsOfAccount, a.RedisClient)
	c.JSON(http.StatusOK, postsOfAccount)
	return
}

func (a *PostsAPI) UpdatePost(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	post := CreatePostStruct(mapBody, accountID)
	post.ID = utils.StringNullable(mapBody["id"])

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
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	mapBody, err := utils.ReadBodyAndReturnMapBody(c.Request.Body)
	if err != nil {
		log.Println(err)
	}

	post := CreatePostStruct(mapBody, accountID)
	post.ID = utils.StringNullable(mapBody["id"])
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

func (a *PostsAPI) SearchPostByAccountFollowing(c *gin.Context) {

	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Token Invalid",
		})
		return
	}

	resCache, err := cache.FindInCache(c.Request, a.RedisClient)
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
			"Message": "Page is not a number",
		})
		return
	}
	postsOfAccount, err := a.Controller.FindPostByAccountFollowingByAccountID(accountID, &page)
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

	cache.InsertCache(c.Request, postsOfAccount, a.RedisClient)
	c.JSON(http.StatusOK, postsOfAccount)
	return
}

func CreatePostStruct(mapBody map[string]interface{}, accountID *string) *entities.Post {

	return &entities.Post{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		Content:   utils.StringNullable(mapBody["content"]),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}
}
