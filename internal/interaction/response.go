package interaction

type InteractionResponse struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	PostID    string `json:"post_id"`
	CommentID string `json:"comment_id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
