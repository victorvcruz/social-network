package controllers

import (
	"social_network_project/controllers/errors"
	"social_network_project/database/repository"
	"social_network_project/entities"
	"social_network_project/entities/response"
)

type CommentsController interface {
	InsertComment(comment *entities.Comment) error
	FindCommentsByAccountID(accountID, postID, commentID *string) ([]interface{}, error)
	UpdateCommentDataByID(comment *entities.Comment) (*response.CommentResponse, error)
	RemoveCommentByID(comment *entities.Comment, accountID *string) (*response.CommentResponse, error)
}

type CommentsControllerStruct struct {
	repositoryComment repository.CommentRepository
	repositoryAccount repository.AccountRepository
	repositoryPost    repository.PostRepository
}

func NewCommentsController() CommentsController {
	return &CommentsControllerStruct{
		repositoryComment: repository.NewComentRepository(),
		repositoryAccount: repository.NewAccountRepository(),
		repositoryPost:    repository.NewPostRepository(),
	}
}

func (c *CommentsControllerStruct) InsertComment(comment *entities.Comment) error {

	existID, err := c.repositoryAccount.ExistsAccountByID(&comment.AccountID)
	if err != nil {
		return err
	}
	if !*existID {
		return &errors.NotFoundAccountIDError{}
	}

	if comment.CommentID != "" {
		existID, err = c.repositoryComment.ExistsCommentByID(&comment.CommentID)
		if err != nil {
			return err
		}
		if !*existID {
			return &errors.NotFoundCommentIDError{}
		}
	}

	err = c.repositoryComment.InsertComment(comment)
	if err != nil {
		return &errors.NotFoundPostIDError{}
	}

	return nil
}

func (c *CommentsControllerStruct) FindCommentsByAccountID(accountID, postID, commentID *string) ([]interface{}, error) {

	existID, err := c.repositoryAccount.ExistsAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	if *postID != "" {
		existID, err = c.repositoryPost.ExistsPostByID(postID)
		if err != nil {
			return nil, err
		}
		if !*existID {
			return nil, &errors.NotFoundPostIDError{}
		}
	}

	if *commentID != "" {
		existID, err = c.repositoryComment.ExistsCommentByID(commentID)
		if err != nil {
			return nil, err
		}
		if !*existID {
			return nil, &errors.NotFoundCommentIDError{}
		}
	}

	return c.repositoryComment.FindCommentsByAccountID(accountID, postID, commentID)
}

func (c *CommentsControllerStruct) UpdateCommentDataByID(comment *entities.Comment) (*response.CommentResponse, error) {

	existID, err := c.repositoryAccount.ExistsAccountByID(&comment.AccountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	err = c.repositoryComment.UpdateCommentDataByID(&comment.ID, comment.Content)
	if err != nil {
		return nil, &errors.NotFoundCommentIDError{}
	}

	postUpdated, err := c.repositoryComment.FindCommentByID(&comment.ID)
	if err != nil {
		return nil, &errors.NotFoundCommentIDError{}
	}

	return postUpdated, nil
}

func (p CommentsControllerStruct) RemoveCommentByID(comment *entities.Comment, accountID *string) (*response.CommentResponse, error) {

	existID, err := p.repositoryAccount.ExistsAccountByID(accountID)

	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	commentToRemoved, err := p.repositoryComment.FindCommentByID(&comment.ID)
	if err != nil {
		return nil, &errors.NotFoundCommentIDError{}
	}

	err = p.repositoryComment.RemoveCommentByID(&comment.ID)
	if err != nil {
		return nil, err
	}

	return commentToRemoved, nil
}
