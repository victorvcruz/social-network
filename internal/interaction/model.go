package interaction

import (
	"database/sql"
)

type Interaction struct {
	ID        string `validate:"required"`
	AccountID string
	PostID    sql.NullString
	CommentID sql.NullString
	Type      InteractionType `validate:"gte=0,lte=1"`
	CreatedAt string
	UpdatedAt string
	Removed   bool
}

func (a *Interaction) ToResponse() *InteractionResponse {
	return &InteractionResponse{
		ID:        a.ID,
		AccountID: a.AccountID,
		PostID:    a.PostID.String,
		CommentID: a.CommentID.String,
		Type:      a.Type.ToString(),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

type InteractionType int

const (
	INTERACTION_TYPE_LIKE InteractionType = iota
	INTERACTION_TYPE_DISLIKE
)

func (i InteractionType) ToString() string {
	return [...]string{"LIKE", "DISLIKE"}[i]
}

func (i InteractionType) EnumIndex() int {
	return int(i)
}

func (i Interaction) ParseString(str string) (InteractionType, bool) {
	var interaction InteractionType
	if str == "LIKE" {
		interaction = INTERACTION_TYPE_LIKE
		return interaction, true
	}
	if str == "DISLIKE" {
		interaction = INTERACTION_TYPE_DISLIKE
		return interaction, true
	}
	return 0, false
}
