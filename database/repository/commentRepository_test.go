package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"testing"
	"time"
)

var c = &entities.Comment{
	ID:        uuid.New().String(),
	AccountID: "f981d822-7efb-4e66-aa84-99f517820ca3",
	PostID:    "0d0bb472-225c-4c8a-9935-a21045c80d87",
	CommentID: "8b607c43-0190-4c8c-9746-4b527d1d2c55",
	Content:   "Lorem ipsum dolor sit amet, consectetuer adipiscing elit.",
	CreatedAt: time.Now().UTC().Format("2006-01-02"),
	UpdatedAt: time.Now().UTC().Format("2006-01-02"),
	Removed:   false,
}

func TestCommentRepositoryStruct_InsertComment(t *testing.T) {
	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		INSERT INTO comment (id, account_id, post_id, comment_id, content, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(c.ID, c.AccountID, c.PostID, c.CommentID, c.Content, c.CreatedAt,
		c.UpdatedAt, c.Removed).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.InsertComment(c)
	assert.Error(t, err)
}
