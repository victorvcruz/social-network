package entities

import (
	"database/sql"
	"social_network_project/entities/response"
)

type Interaction struct {
	ID        string `validate:"required"`
	AccountID string
	PostID    sql.NullString
	CommentID sql.NullString
	Type      response.InteractionType `validate:"gte=0,lte=1"`
	CreatedAt string
	UpdatedAt string
	Removed   bool
}

func (a *Interaction) ToResponse() *response.InteractionResponse {
	return &response.InteractionResponse{
		ID:        a.ID,
		AccountID: a.AccountID,
		PostID:    a.PostID.String,
		CommentID: a.CommentID.String,
		Type:      a.Type,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
