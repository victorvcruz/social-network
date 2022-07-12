package response

import (
	"strings"
)

type InteractionResponse struct {
	ID        string
	AccountID string
	PostID    string
	CommentID string
	Type      InteractionType
	CreatedAt string
	UpdatedAt string
}

type InteractionType int

const (
	INTERACTION_TYPE_LIKED    InteractionType = 0
	INTERACTION_TYPE_DISLIKED InteractionType = 1
)

var (
	interactionsMap = map[string]InteractionType{
		"like":    INTERACTION_TYPE_LIKED,
		"dislike": INTERACTION_TYPE_DISLIKED,
	}
)

func (i InteractionType) EnumIndex() int {
	return int(i)
}

func ParseString(str string) (InteractionType, bool) {
	c, ok := interactionsMap[strings.ToLower(str)]
	return c, ok
}
