package handlers

import (
	"github.com/stretchr/testify/assert"
	"social_network_project/internal/comment"
	"social_network_project/internal/utils"
	"testing"
	"time"
)

func TestCommentsAPI_CreateCommentStruct(t *testing.T) {

	var mapBody = map[string]interface{}{
		"content": "Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
	}

	ID := "74657673-2405-4d47-afc3-5268ed1370c9"
	accountID := "6ed3f907-26dc-411d-9b85-569b66423d6f"
	postID := "25cc459c-f1cd-42e0-8b7b-b28853e3a68a"
	commentID := utils.NewNullString("f97ff1b9-9ee7-4ec5-addf-a9983aefd571")

	commentExpected := &comment.Comment{
		ID:        ID,
		AccountID: accountID,
		PostID:    postID,
		CommentID: commentID,
		Content:   utils.StringNullable(mapBody["content"]),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

	comment := CreateCommentStruct(mapBody, &accountID, &postID, commentID)
	comment.ID = ID

	assert.Equal(t, commentExpected, comment)
}
