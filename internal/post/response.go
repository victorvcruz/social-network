package post

type PostResponse struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Like      int `json:"like"`
	Dislike   int `json:"dislike"`
}
