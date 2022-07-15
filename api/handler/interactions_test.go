package handler

import (
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"social_network_project/utils"
	"testing"
	"time"
)

func TestInteractionsAPI_CreateInteractionStruct(t *testing.T) {

	var mapBody = map[string]interface{}{
		"post_id": "baa4570c-1904-454d-87f7-cdd9033b94af",
		"type":    "LIKE",
	}

	ID := "74657673-2405-4d47-afc3-5268ed1370c9"
	accountID := "6ed3f907-26dc-411d-9b85-569b66423d6f"

	interactionExpected := &entities.Interaction{
		ID:        ID,
		AccountID: accountID,
		PostID:    utils.NewNullString(utils.StringNullable(mapBody["post_id"])),
		CommentID: utils.NewNullString(utils.StringNullable(mapBody["comment_id"])),
		Type:      0,
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

	interaction := CreateInteractionStruct(mapBody, &accountID)
	interaction.ID = ID

	assert.Equal(t, interactionExpected, interaction)
}
