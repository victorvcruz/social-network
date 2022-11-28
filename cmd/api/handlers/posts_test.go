package handlers

import (
	"github.com/stretchr/testify/assert"
	"social_network_project/internal/post"
	"testing"
	"time"
)

func TestPostsAPI_CreatePostStruct(t *testing.T) {

	var mapBody = map[string]interface{}{
		"content": "Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
	}

	accountID := "6ed3f907-26dc-411d-9b85-569b66423d6f"
	postID := "25cc459c-f1cd-42e0-8b7b-b28853e3a68a"

	postExpected := &post.Post{
		ID:        postID,
		AccountID: accountID,
		Content:   mapBody["content"].(string),
		CreatedAt: time.Now().UTC().Format("2006-01-02"),
		UpdatedAt: time.Now().UTC().Format("2006-01-02"),
		Removed:   false,
	}

	post := CreatePostStruct(mapBody, &accountID)
	post.ID = postID

	assert.Equal(t, postExpected, post)
}
