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
	interaction2 "social_network_project/internal/interaction"
	"social_network_project/internal/interaction/service"
	"social_network_project/internal/utils"
	"social_network_project/internal/utils/errors"
	"social_network_project/internal/utils/validate"
	"time"
)

type IntercationsHandlerClient interface {
	CreateInteraction(c *gin.Context)
	UpdateInteraction(c *gin.Context)
	DeleteInteraction(c *gin.Context)
}

type InteractionsHandler struct {
	Controller service.InteractionsServiceClient
	Validate   *validator.Validate
}

func RegisterInteractionsHandlers(interactionsController service.InteractionsServiceClient) IntercationsHandlerClient {
	return &InteractionsHandler{
		Controller: interactionsController,
		Validate:   validator.New(),
	}
}

func (i InteractionsHandler) CreateInteraction(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request interaction2.InteractionRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	if request.PostId == "" && request.CommentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Add post_id or comment_id",
		})
		return
	}

	interaction := i.fillFields(request, &accountID)

	mapper := make(map[string]interface{})
	err = i.Validate.Struct(interaction)
	if err != nil {
		mapper["errors"] = validate.RequestInteractionValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	err = i.Controller.InsertInteraction(interaction)
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
		case *errors.ConflictAlreadyWriteError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, interaction.ToResponse())
	return
}

func (i *InteractionsHandler) UpdateInteraction(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request interaction2.InteractionRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	interaction := i.fillFields(request, &accountID)
	interaction.ID = utils.StringNullable(request.Id)

	mapper := make(map[string]interface{})
	err = i.Validate.Struct(interaction)
	if err != nil {
		mapper["errors"] = validate.RequestInteractionValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	interactionUpdated, err := i.Controller.UpdateInteractonDataByID(interaction)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundInteractionIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, interactionUpdated)
	return
}

func (i *InteractionsHandler) DeleteInteraction(c *gin.Context) {
	accountID, err := utils.DecodeTokenAndReturnID(c.Request.Header.Get(os.Getenv("JWT_TOKEN_HEADER")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token Invalid",
		})
		return
	}

	var request interaction2.InteractionRequest

	body, err := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unprocessable Entity",
		})
		return
	}

	interaction := i.fillFields(request, &accountID)
	interaction.ID = utils.StringNullable(request.Id)
	interaction.Type = 0

	mapper := make(map[string]interface{})
	err = i.Validate.Struct(interaction)
	if err != nil {
		mapper["errors"] = validate.RequestInteractionValidate(err)
		c.JSON(http.StatusBadRequest, mapper)
		return
	}

	interactionRemoved, err := i.Controller.RemoveInteractionByID(interaction)
	if err != nil {
		switch e := err.(type) {
		case *errors.UnauthorizedAccountIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		case *errors.NotFoundInteractionIDError:
			log.Println(e)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		default:
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, interactionRemoved)
	return
}

func (i *InteractionsHandler) fillFields(req interaction2.InteractionRequest, accountID *string) *interaction2.Interaction {
	interaction := &interaction2.Interaction{
		ID:        uuid.New().String(),
		AccountID: *accountID,
		PostID:    utils.NewNullString(utils.StringNullable(req.PostId)),
		CommentID: utils.NewNullString(utils.StringNullable(req.CommentId)),
		Type:      0,
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}
	typeInteraction, ok := interaction.ParseString(utils.StringNullable(req.Type))
	if !ok {
		typeInteraction = 400
	}
	interaction.Type = typeInteraction
	return interaction
}
