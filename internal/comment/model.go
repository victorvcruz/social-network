package comment

import (
	"database/sql"
)

type Comment struct {
	ID        string `validate:"required"`
	AccountID string
	PostID    string
	CommentID sql.NullString
	Content   string `validate:"required"`
	CreatedAt string
	UpdatedAt string
	Removed   bool
	Like      int
	Dislike   int
}

func (a *Comment) ToResponse() *CommentResponse {
	return &CommentResponse{
		ID:        a.ID,
		AccountID: a.AccountID,
		PostID:    a.PostID,
		CommentID: a.CommentID.String,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Like:      a.Like,
		Dislike:   a.Dislike,
	}
}
