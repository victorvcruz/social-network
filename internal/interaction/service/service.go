package service

import (
	"social_network_project/internal/account"
	"social_network_project/internal/comment"
	"social_network_project/internal/interaction"
	"social_network_project/internal/notification"
	"social_network_project/internal/notification/service"
	"social_network_project/internal/utils/errors"
)

type InteractionsServiceClient interface {
	InsertInteraction(interaction *interaction.Interaction) error
	UpdateInteractonDataByID(interaction *interaction.Interaction) (*interaction.InteractionResponse, error)
	RemoveInteractionByID(interaction *interaction.Interaction) (*interaction.InteractionResponse, error)
}

type InteractionsService struct {
	repositoryAccount     account.AccountRepository
	repositoryComment     comment.CommentRepository
	repositoryInteraction interaction.InteractionRepository
	rabbitControl         service.NotificationServiceClient
}

func NewInteractionsService(_repositoryAccount account.AccountRepository, _repositoryComment comment.CommentRepository, _repositoryInteraction interaction.InteractionRepository, rabbitmq service.NotificationServiceClient) InteractionsServiceClient {
	return &InteractionsService{
		repositoryAccount:     _repositoryAccount,
		repositoryComment:     _repositoryComment,
		repositoryInteraction: _repositoryInteraction,
		rabbitControl:         rabbitmq,
	}
}

func (i InteractionsService) InsertInteraction(interaction *interaction.Interaction) error {

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

	i.rabbitControl.SendMessage(notification.CreateNotificationJson("Interaction", interaction.ID))
	return nil
}

func (i InteractionsService) UpdateInteractonDataByID(interaction *interaction.Interaction) (*interaction.InteractionResponse, error) {

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

func (i InteractionsService) RemoveInteractionByID(interaction *interaction.Interaction) (*interaction.InteractionResponse, error) {

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
