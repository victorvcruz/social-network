package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"social_network_project/controllers"
	"social_network_project/controllers/errors"
	"social_network_project/controllers/validate"
	"social_network_project/entities"
	"social_network_project/utils"
	"time"
)

type InteractionsAPI struct {
	Controller controllers.InteractionsController
	Validate   *validator.Validate
}

func RegisterInteractionsHandlers(handler *gin.Engine, interactionsController controllers.InteractionsController) {
	ac := &InteractionsAPI{
		Controller: interactionsController,
		Validate:   validator.New(),
	}

	handler.POST("/interaction", ac.CreateInteraction)
	handler.PUT("/interaction", ac.UpdateInteraction)
	handler.DELETE("/interaction", ac.DeleteInteraction)
}

func (a InteractionsAPI) CreateInteraction(c *gin.Context) {
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

	if mapBody["post_id"] == nil && mapBody["comment_id"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Add post_id or comment_id",
		})
		return
	}

	interaction := CreateInteractionStruct(mapBody, accountID)

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(interaction)
	if err != nil {
		mapper["errors"] = validate.RequestInteractionValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	err = a.Controller.InsertInteraction(interaction)
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
		case *errors.ConflictAlreadyWriteError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, interaction.ToResponse())
	return
}

func (a InteractionsAPI) UpdateInteraction(c *gin.Context) {
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

	interaction := CreateInteractionStruct(mapBody, accountID)
	interaction.ID = utils.StringNullable(mapBody["id"])

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(interaction)
	if err != nil {
		mapper["errors"] = validate.RequestInteractionValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	interactionUpdated, err := a.Controller.UpdateInteractonDataByID(interaction)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundInteractionIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, interactionUpdated)
	return
}

func (a InteractionsAPI) DeleteInteraction(c *gin.Context) {
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

	interaction := CreateInteractionStruct(mapBody, accountID)
	interaction.ID = utils.StringNullable(mapBody["id"])
	interaction.Type = 0

	mapper := make(map[string]interface{})
	err = a.Validate.Struct(interaction)
	if err != nil {
		mapper["errors"] = validate.RequestInteractionValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	interactionRemoved, err := a.Controller.RemoveInteractionByID(interaction)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		case *errors.NotFoundInteractionIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"Message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, interactionRemoved)
	return
}

func CreateInteractionStruct(mapBody map[string]interface{}, accountID *string) *entities.Interaction {

	interaction, ok := entities.ParseString(utils.StringNullable(mapBody["type"]))
	if !ok {
		interaction = 400
	}

	return &entities.Interaction{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		PostID:    utils.NewNullString(utils.StringNullable(mapBody["post_id"])),
		CommentID: utils.NewNullString(utils.StringNullable(mapBody["comment_id"])),
		Type:      interaction,
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}
}
