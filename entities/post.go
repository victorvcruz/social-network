package entities

import "social_network_project/entities/response"

type Post struct {
	ID        string `validate:"required"`
	AccountID string
	Content   string `validate:"required"`
	CreatedAt string
	UpdatedAt string
	Removed   bool
}

func (a *Post) ToResponse() response.PostResponse {
	return response.PostResponse{
		ID:        a.ID,
		AccountID: a.AccountID,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
