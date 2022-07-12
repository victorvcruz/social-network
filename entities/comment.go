package entities

import (
	"database/sql"
	"social_network_project/entities/response"
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
}

func (a *Comment) ToResponse() *response.CommentResponse {
	return &response.CommentResponse{
		ID:        a.ID,
		AccountID: a.AccountID,
		PostID:    a.PostID,
		CommentID: a.CommentID.String,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
