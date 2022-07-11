package response

type CommentResponse struct {
	ID        string
	AccountID string
	PostID    string
	CommentID string
	Content   string
	CreatedAt string
	UpdatedAt string
}
