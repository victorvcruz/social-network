package controllers

import (
	"social_network_project/controllers/errors"
	"social_network_project/database/repository"
	"social_network_project/entities"
	"social_network_project/entities/response"
)

type InteractionsController interface {
	InsertInteraction(interaction *entities.Interaction) error
	UpdateInteractonDataByID(interaction *entities.Interaction) (*response.InteractionResponse, error)
	RemoveInteractionByID(interaction *entities.Interaction) (*response.InteractionResponse, error)
}

type InteractionsControllerStruct struct {
	repositoryAccount     repository.AccountRepository
	repositoryPost        repository.PostRepository
	repositoryComment     repository.CommentRepository
	repositoryInteraction repository.InteractionRepository
}

func NewInteractionsController() InteractionsController {
	return &InteractionsControllerStruct{
		repositoryAccount:     repository.NewAccountRepository(),
		repositoryPost:        repository.NewPostRepository(),
		repositoryComment:     repository.NewComentRepository(),
		repositoryInteraction: repository.NewInteractionRepository(),
	}
}

func (i InteractionsControllerStruct) InsertInteraction(interaction *entities.Interaction) error {

	existID, err := i.repositoryAccount.ExistsAccountByID(&interaction.AccountID)
	if err != nil {
		return err
	}
	if !*existID {
		return &errors.NotFoundAccountIDError{}
	}

	if interaction.CommentID.String != "" {
		existID, err = i.repositoryComment.ExistsCommentByID(&interaction.CommentID.String)
		if err != nil {
			return err
		}
		if !*existID {
			return &errors.NotFoundCommentIDError{}
		}
		existID, err = i.repositoryInteraction.ExistsInteractionByCommentIDAndAccountID(&interaction.CommentID.String, &interaction.AccountID)
		if err != nil {
			return err
		}
		if *existID {
			return &errors.ConflictAlreadyWriteError{}
		}
	}

	if interaction.PostID.String != "" {
		existID, err = i.repositoryInteraction.ExistsInteractionByPostIDAndAccountID(&interaction.PostID.String, &interaction.AccountID)
		if err != nil {
			return err
		}
		if *existID {
			return &errors.ConflictAlreadyWriteError{}
		}
	}

	err = i.repositoryInteraction.InsertInteraction(interaction)
	if err != nil {
		return &errors.NotFoundPostIDError{}
	}

	return nil
}

func (i InteractionsControllerStruct) UpdateInteractonDataByID(interaction *entities.Interaction) (*response.InteractionResponse, error) {

	exist, err := i.repositoryInteraction.ExistsInteractionByID(&interaction.ID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.NotFoundInteractionIDError{}
	}

	exist, err = i.repositoryInteraction.ExistsInteractionByInteractionIDAndAccountID(&interaction.ID, &interaction.AccountID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.UnauthorizedAccountIDError{}
	}

	err = i.repositoryInteraction.UpdateInteractonDataByID(&interaction.ID, &interaction.AccountID, &interaction.Type)
	if err != nil {
		return nil, err
	}

	interactionUpdated, err := i.repositoryInteraction.FindInteractionByID(&interaction.ID)
	if err != nil {
		return nil, &errors.NotFoundInteractionIDError{}
	}

	return interactionUpdated.ToResponse(), nil
}

func (i InteractionsControllerStruct) RemoveInteractionByID(interaction *entities.Interaction) (*response.InteractionResponse, error) {

	interactionToRemove, err := i.repositoryInteraction.FindInteractionByID(&interaction.ID)
	if err != nil {
		return nil, &errors.NotFoundInteractionIDError{}
	}

	existID, err := i.repositoryInteraction.ExistsInteractionByInteractionIDAndAccountID(&interaction.ID, &interaction.AccountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.UnauthorizedAccountIDError{}
	}

	err = i.repositoryInteraction.RemoveInteractionByID(&interaction.ID, &interaction.AccountID)
	if err != nil {
		return nil, err
	}

	return interactionToRemove.ToResponse(), nil
}
