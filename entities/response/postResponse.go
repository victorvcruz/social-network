package response

type PostResponse struct {
	ID        string
	AccountID string
	Content   string
	CreatedAt string
	UpdatedAt string
	Like      int
	Dislike   int
}
