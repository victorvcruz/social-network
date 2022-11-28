package comment

type CommentRequest struct {
	Id      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}