package service

import (
	"social_network_project/internal/account"
	"social_network_project/internal/notification"
	"social_network_project/internal/notification/service"
	"social_network_project/internal/post"
	"social_network_project/internal/utils/errors"
)

type PostsServiceClient interface {
	InsertPost(post *post.Post) error
	FindPostsByAccountID(accountID, idToGet, page *string) ([]interface{}, error)
	UpdatePostDataByID(post *post.Post) (*post.PostResponse, error)
	RemovePostByID(post *post.Post) (*post.PostResponse, error)
	FindPostByAccountFollowingByAccountID(accountID *string, page *string) ([]interface{}, error)
}

type PostsService struct {
	repositoryPost    post.PostRepository
	repositoryAccount account.AccountRepository
	rabbitControl     service.NotificationServiceClient
}

func NewPostsService(_repositoryPost post.PostRepository, _repositoryAccount account.AccountRepository, rabbitmq service.NotificationServiceClient) PostsServiceClient {
	return &PostsService{
		repositoryPost:   _repositoryPost,
		repositoryAccount: _repositoryAccount,
		rabbitControl:     rabbitmq,
	}
}

func (p PostsService) InsertPost(post *post.Post) error {

	err := p.repositoryPost.InsertPost(post)
	if err != nil {
		return &errors.NotFoundAccountIDError{}
	}

	p.rabbitControl.SendMessage(notification.CreateNotificationJson("Post", post.AccountID))
	return nil
}

func (p PostsService) FindPostsByAccountID(accountID, idToGet, page *string) ([]interface{}, error) {

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

func (p PostsService) UpdatePostDataByID(post *post.Post) (*post.PostResponse, error) {

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

func (p PostsService) RemovePostByID(post *post.Post) (*post.PostResponse, error) {

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

func (p PostsService) FindPostByAccountFollowingByAccountID(accountID *string, page *string) ([]interface{}, error) {

	existID, err := p.repositoryAccount.ExistsAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	return p.repositoryPost.FindPostByAccountFollowingByAccountID(accountID, page)
}
