package controllers

import (
	"database/sql"
	"social_network_project/controllers/errors"
	"social_network_project/database/repository"
	"social_network_project/entities"
	"social_network_project/entities/response"
)

type CommentsController interface {
	InsertComment(comment *entities.Comment) error
	FindCommentsByAccountID(accountID, idToGet, postID, commentID, page *string) ([]interface{}, error)
	UpdateCommentDataByID(comment *entities.Comment) (*response.CommentResponse, error)
	RemoveCommentByID(comment *entities.Comment) (*response.CommentResponse, error)
}

type CommentsControllerStruct struct {
	repositoryComment repository.CommentRepository
	repositoryAccount repository.AccountRepository
	repositoryPost    repository.PostRepository
}

func NewCommentsController(postgresDB *sql.DB) CommentsController {
	return &CommentsControllerStruct{
		repositoryComment: repository.NewComentRepository(postgresDB),
		repositoryAccount: repository.NewAccountRepository(postgresDB),
		repositoryPost:    repository.NewPostRepository(postgresDB),
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

	if comment.CommentID.String != "" {
		existID, err = c.repositoryComment.ExistsCommentByID(&comment.CommentID.String)
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

func (c *CommentsControllerStruct) FindCommentsByAccountID(accountID, idToGet, postID, commentID, page *string) ([]interface{}, error) {

	existID, err := c.repositoryAccount.ExistsAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.NotFoundAccountIDError{}
	}

	if *idToGet != "" {
		existID, err := c.repositoryAccount.ExistsAccountByID(idToGet)
		if err != nil {
			return nil, err
		}
		if !*existID {
			return nil, &errors.NotFoundAccountIDError{}
		}
		accountID = idToGet
	}

	if *postID != "" {
		existID, err = c.repositoryPost.ExistsPostByID(postID)
		if err != nil {
			return nil, err
		}
		if !*existID {
			return nil, &errors.NotFoundPostIDError{}
		}
		return c.repositoryComment.FindCommentsByPostOrCommentID(postID, commentID, page)

	}

	if *commentID != "" {
		existID, err = c.repositoryComment.ExistsCommentByID(commentID)
		if err != nil {
			return nil, err
		}
		if !*existID {
			return nil, &errors.NotFoundCommentIDError{}
		}
		return c.repositoryComment.FindCommentsByPostOrCommentID(postID, commentID, page)
	}

	return c.repositoryComment.FindCommentsByAccountID(accountID, page)
}

func (c *CommentsControllerStruct) UpdateCommentDataByID(comment *entities.Comment) (*response.CommentResponse, error) {

	exist, err := c.repositoryComment.ExistsCommentByID(&comment.ID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.NotFoundCommentIDError{}
	}

	exist, err = c.repositoryComment.ExistsCommentByCommentIDAndAccountID(&comment.ID, &comment.AccountID)
	if err != nil {
		return nil, err
	}
	if !*exist {
		return nil, &errors.UnauthorizedAccountIDError{}
	}

	err = c.repositoryComment.UpdateCommentDataByID(&comment.ID, &comment.AccountID, &comment.Content)
	if err != nil {
		return nil, err
	}

	postUpdated, err := c.repositoryComment.FindCommentByID(&comment.ID)
	if err != nil {
		return nil, &errors.NotFoundCommentIDError{}
	}

	return postUpdated.ToResponse(), nil
}

func (p CommentsControllerStruct) RemoveCommentByID(comment *entities.Comment) (*response.CommentResponse, error) {

	commentToRemoved, err := p.repositoryComment.FindCommentByID(&comment.ID)
	if err != nil {
		return nil, &errors.NotFoundCommentIDError{}
	}

	existID, err := p.repositoryComment.ExistsCommentByCommentIDAndAccountID(&comment.ID, &comment.AccountID)
	if err != nil {
		return nil, err
	}
	if !*existID {
		return nil, &errors.UnauthorizedAccountIDError{}
	}

	err = p.repositoryComment.RemoveCommentByID(&comment.ID, &comment.AccountID)
	if err != nil {
		return nil, err
	}

	return commentToRemoved.ToResponse(), nil
}
