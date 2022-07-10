package entities

import "social_network_project/entities/response"

type Post struct {
	ID        string
	AccountID string
	Content   string
	CreatedAt string
	UpdatedAt string
	Removed   bool
}

func (a *Post) ToResponse() response.PostResponse {
	return response.PostResponse{
		ID:        a.ID,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
