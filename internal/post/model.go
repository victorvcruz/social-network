package post

type Post struct {
	ID        string `validate:"required"`
	AccountID string
	Content   string `validate:"required"`
	CreatedAt string
	UpdatedAt string
	Removed   bool
	Like      int
	Dislike   int
}

func (a *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:        a.ID,
		AccountID: a.AccountID,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Like:      a.Like,
		Dislike:   a.Dislike,
	}
}
