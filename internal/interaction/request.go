package interaction

type InteractionRequest struct {
	Id        string `json:"id,omitempty"`
	PostId    string `json:"post_id"`
	CommentId string `json:"comment_id,omitempty"`
	Type      string `json:"type,omitempty"`
}