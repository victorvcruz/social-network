package entities

import "social_network_project/entities/response"

type Comment struct {
	ID        string
	AccountID string
	PostID    string
	CommentID string
	Content   string
	CreatedAt string
	UpdatedAt string
	Removed   bool
}

func (a *Comment) ToResponse() response.CommentResponse {
	return response.CommentResponse{
		ID:        a.ID,
		PostID:    a.PostID,
		CommentID: a.CommentID,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
