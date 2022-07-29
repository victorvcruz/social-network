package controllers

import (
	"database/sql"
	"social_network_project/controllers/errors"
	"social_network_project/database/repository"
	"social_network_project/entities"
	"social_network_project/entities/response"
	message_broker "social_network_project/message-broker"
)

type PostsController interface {
	InsertPost(post *entities.Post) error
	FindPostsByAccountID(accountID, idToGet, page *string) ([]interface{}, error)
	UpdatePostDataByID(post *entities.Post) (*response.PostResponse, error)
	RemovePostByID(post *entities.Post) (*response.PostResponse, error)
	FindPostByAccountFollowingByAccountID(accountID *string, page *string) ([]interface{}, error)
}

type PostsControllerStruct struct {
	repositoryPost    repository.PostRepository
	repositoryAccount repository.AccountRepository
	rabbitControl     *message_broker.NotificationControl
}

func NewPostsController(postgresDB *sql.DB, rabbitmq *message_broker.NotificationControl) PostsController {
	return &PostsControllerStruct{
		repositoryPost:    repository.NewPostRepository(postgresDB),
		repositoryAccount: repository.NewAccountRepository(postgresDB),
		rabbitControl:     rabbitmq,
	}
}

func (p PostsControllerStruct) InsertPost(post *entities.Post) error {

	err := p.repositoryPost.InsertPost(post)
	if err != nil {
		return &errors.NotFoundAccountIDError{}
	}

	p.rabbitControl.SendMessage(message_broker.CreateNotificationJson("Post", post.AccountID))
	return nil
}

func (p PostsControllerStruct) FindPostsByAccountID(accountID, idToGet, page *string) ([]interface{}, error) {

	existID, err := p.repositoryAccount.ExistsAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	if *idToGet != "" {
		existID, err := p.repositoryAccount.ExistsAccountByID(idToGet)
		if err != nil {
			return nil, err
		}
		if !*existID {
			return nil, &errors.NotFoundAccountIDError{}
		}

		return p.repositoryPost.FindPostsByAccountID(idToGet, page)
	}

	return p.repositoryPost.FindPostsByAccountID(accountID, page)
}

func (p PostsControllerStruct) UpdatePostDataByID(post *entities.Post) (*response.PostResponse, error) {

	exist, err := p.repositoryPost.ExistsPostByID(&post.ID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.NotFoundPostIDError{}
	}

	exist, err = p.repositoryPost.ExistsPostByPostIDAndAccountID(&post.ID, &post.AccountID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.UnauthorizedAccountIDError{}
	}

	err = p.repositoryPost.UpdatePostDataByID(&post.ID, &post.AccountID, &post.Content)
	if err != nil {
		return nil, err
	}

	postUpdated, err := p.repositoryPost.FindPostByID(&post.ID)
	if err != nil {
		return nil, &errors.NotFoundPostIDError{}
	}

	return postUpdated, nil
}

func (p PostsControllerStruct) RemovePostByID(post *entities.Post) (*response.PostResponse, error) {

	postToRemoved, err := p.repositoryPost.FindPostByID(&post.ID)
	if err != nil {
		return nil, &errors.NotFoundPostIDError{}
	}

	exist, err := p.repositoryPost.ExistsPostByPostIDAndAccountID(&post.ID, &post.AccountID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.UnauthorizedAccountIDError{}
	}

	err = p.repositoryPost.RemovePostByID(&post.ID, &post.AccountID)
	if err != nil {
		return nil, err
	}

	return postToRemoved, nil
}

func (p PostsControllerStruct) FindPostByAccountFollowingByAccountID(accountID *string, page *string) ([]interface{}, error) {

	existID, err := p.repositoryAccount.ExistsAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return p.repositoryPost.FindPostByAccountFollowingByAccountID(accountID, page)
}
