package controllers

import (
	"social_network_project/controllers/errors"
	"social_network_project/database/repository"
	"social_network_project/entities"
	"social_network_project/entities/response"
)

type PostsController interface {
	InsertPost(post *entities.Post) error
	FindPostsByAccountID(id *string) ([]interface{}, error)
	UpdatePostDataByID(post *entities.Post) (*response.PostResponse, error)
	RemovePostByID(post *entities.Post, accountID *string) (*response.PostResponse, error)
}

type PostsControllerStruct struct {
	repositoryPost    repository.PostRepository
	repositoryAccount repository.AccountRepository
}

func NewPostsController() PostsController {
	return &PostsControllerStruct{
		repositoryPost:    repository.NewPostRepository(),
		repositoryAccount: repository.NewAccountRepository(),
	}
}

func (p PostsControllerStruct) InsertPost(post *entities.Post) error {

	err := p.repositoryPost.InsertPost(post)
	if err != nil {
		return &errors.NotFoundAccountIDError{}
	}

	return nil
}

func (p PostsControllerStruct) FindPostsByAccountID(id *string) ([]interface{}, error) {

	existID, err := p.repositoryAccount.ExistsAccountByID(id)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return p.repositoryPost.FindPostsByAccountID(id)
}

func (p PostsControllerStruct) UpdatePostDataByID(post *entities.Post) (*response.PostResponse, error) {

	existID, err := p.repositoryAccount.ExistsAccountByID(&post.AccountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	err = p.repositoryPost.UpdatePostDataByID(&post.ID, post.Content)
	if err != nil {
		return nil, &errors.NotFoundPostIDError{}
	}

	postUpdated, err := p.repositoryPost.FindPostByID(&post.ID)
	if err != nil {
		return nil, &errors.NotFoundPostIDError{}
	}

	return postUpdated, nil
}

func (p PostsControllerStruct) RemovePostByID(post *entities.Post, accountID *string) (*response.PostResponse, error) {

	existID, err := p.repositoryAccount.ExistsAccountByID(accountID)

	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	postToRemoved, err := p.repositoryPost.FindPostByID(&post.ID)
	if err != nil {
		return nil, &errors.NotFoundPostIDError{}
	}

	err = p.repositoryPost.RemovePostByID(&post.ID)
	if err != nil {
		return nil, err
	}

	return postToRemoved, nil
}
