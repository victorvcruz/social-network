package post

type PostRequest struct {
	Id      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}