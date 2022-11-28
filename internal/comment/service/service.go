package service

import (
	"social_network_project/internal/account"
	"social_network_project/internal/comment"
	"social_network_project/internal/notification"
	"social_network_project/internal/notification/service"
	"social_network_project/internal/post"
	"social_network_project/internal/utils/errors"
)

type CommentsServiceClient interface {
	InsertComment(comment *comment.Comment) error
	FindCommentsByAccountID(accountID, idToGet, postID, commentID, page *string) ([]interface{}, error)
	UpdateCommentDataByID(comment *comment.Comment) (*comment.CommentResponse, error)
	RemoveCommentByID(comment *comment.Comment) (*comment.CommentResponse, error)
}

type CommentsService struct {
	repositoryComment comment.CommentRepository
	repositoryAccount account.AccountRepository
	repositoryPost post.PostRepository
	rabbitControl  service.NotificationServiceClient
}

func NewCommentsService(_repositoryComment comment.CommentRepository, _repositoryAccount account.AccountRepository, _repositoryPost post.PostRepository,rabbitmq service.NotificationServiceClient) CommentsServiceClient {
	return &CommentsService{
		repositoryComment: _repositoryComment,
		repositoryAccount: _repositoryAccount,
		repositoryPost:    _repositoryPost,
		rabbitControl:     rabbitmq,
	}
}

func (c *CommentsService) InsertComment(comment *comment.Comment) error {

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

	c.rabbitControl.SendMessage(notification.CreateNotificationJson("Comment", comment.ID))
	return nil
}

func (c *CommentsService) FindCommentsByAccountID(accountID, idToGet, postID, commentID, page *string) ([]interface{}, error) {

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

func (c *CommentsService) UpdateCommentDataByID(comment *comment.Comment) (*comment.CommentResponse, error) {

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

func (p CommentsService) RemoveCommentByID(comment *comment.Comment) (*comment.CommentResponse, error) {

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
