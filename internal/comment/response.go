package comment

type CommentResponse struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	PostID    string `json:"post_id"`
	CommentID string `json:"comment_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Like      int `json:"like"`
	Dislike   int `json:"dislike"`
}
