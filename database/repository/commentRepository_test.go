package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"social_network_project/utils"
	"testing"
	"time"
)

var c = &entities.Comment{
	ID:        uuid.New().String(),
	AccountID: "f981d822-7efb-4e66-aa84-99f517820ca3",
	PostID:    "0d0bb472-225c-4c8a-9935-a21045c80d87",
	CommentID: utils.NewNullString("8b607c43-0190-4c8c-9746-4b527d1d2c55"),
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

func TestCommentRepositoryStruct_ExistsCommentByID(t *testing.T) {
	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM comment
		WHERE id = $1
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(c.ID)

	mock.ExpectQuery(query).WithArgs(c.ID).WillReturnRows(rows)

	exist, err := repo.ExistsCommentByID(&c.ID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestCommentRepositoryStruct_FindCommentsByAccountID(t *testing.T) {

	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()
	page := "1"
	query :=
		`SELECT comment.id, comment.account_id, comment.post_id, comment.comment_id, comment.content, comment.created_at, comment.updated_at,
		(
			SELECT count(1) FROM interaction i WHERE i.comment_id = comment.id AND i."type" = 'LIKE'
	) AS like,
		(
			SELECT count(1) FROM interaction i WHERE i.comment_id = comment.id AND i."type" = 'DISLIKE'
	) AS dislike
	FROM comment
	WHERE comment.account_id = $1
	AND comment.removed = false
	Order By comment.created_at
	OFFSET ($2 - 1) * 10
	FETCH NEXT 10 ROWS ONLY;`

	rows := sqlmock.NewRows([]string{"id", "account_id", "post_id", "comment_id", "content", "created_at", "updated_at"}).
		AddRow(c.ID, c.AccountID, c.PostID, c.CommentID, c.Content, p.CreatedAt, p.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(c.ID).WillReturnRows(rows)

	account, err := repo.FindCommentsByAccountID(&c.AccountID, &page)
	assert.Empty(t, account)
	assert.Error(t, err)
}

func TestCommentRepositoryStruct_FindCommentsByPostOrCommentID(t *testing.T) {

	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()
	page := "1"
	query :=
		`SELECT comment.id, comment.account_id, comment.post_id, comment.comment_id, comment.content, comment.created_at, comment.updated_at,
		(
			SELECT count(1) FROM interaction i WHERE i.comment_id = comment.id AND i."type" = 'LIKE'
	) AS like,
		(
			SELECT count(1) FROM interaction i WHERE i.comment_id = comment.id AND i."type" = 'DISLIKE'
	) AS dislike
	FROM comment
	WHERE ` + "comment.comment_id = $1 " +
			`AND comment.removed = false
	Order By comment.created_at
	OFFSET ($2 - 1) * 10
	FETCH NEXT 10 ROWS ONLY;`

	rows := sqlmock.NewRows([]string{"id", "account_id", "post_id", "comment_id", "content", "created_at", "updated_at"}).
		AddRow(c.ID, c.AccountID, c.PostID, c.CommentID, c.Content, p.CreatedAt, p.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(c.ID).WillReturnRows(rows)
	t.Run("commentID nil", func(t *testing.T) {
		postID := "09c17021-73a7-43e0-a63b-a237d1f3b85e"
		commentID := ""
		account, err := repo.FindCommentsByPostOrCommentID(&postID, &commentID, &page)
		assert.Empty(t, account)
		assert.Error(t, err)
	})
	t.Run("postID nil", func(t *testing.T) {
		postID := ""
		commentID := "09c17021-73a7-43e0-a63b-a237d1f3b85e"
		account, err := repo.FindCommentsByPostOrCommentID(&postID, &commentID, &page)
		assert.Empty(t, account)
		assert.Error(t, err)
	})

}

func TestCommentRepositoryStruct_UpdateCommentDataByID(t *testing.T) {

	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		UPDATE comment 
		SET content = $1, updated_at = $2
		WHERE id = $3
		AND account_id = $2
		AND removed = false`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(c.Content, c.UpdatedAt, c.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateCommentDataByID(&c.ID, &c.AccountID, &c.Content)
	assert.Error(t, err)
}

func TestCommentRepositoryStruct_FindCommentByID(t *testing.T) {
	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id, account_id, content, created_at, updated_at
		FROM comment
		WHERE id = $1
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id", "account_id", "content", "created_at", "updated_at"}).
		AddRow(p.ID, p.AccountID, p.Content, p.CreatedAt, p.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(p.ID).WillReturnRows(rows)

	account, err := repo.FindCommentByID(&p.ID)
	assert.Empty(t, account)
	assert.Error(t, err)
}

func TestCommentRepositoryStruct_RemoveCommentByID(t *testing.T) {

	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		UPDATE comment 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(c.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.RemoveCommentByID(&c.ID, &c.AccountID)
	assert.Error(t, err)
}

func TestCommentRepositoryStruct_ExistsCommentByCommentIDAndAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM comment
		WHERE id = $1
		AND account_id = $2
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(c.ID)

	mock.ExpectQuery(query).WithArgs(c.ID).WillReturnRows(rows)

	exist, err := repo.ExistsCommentByCommentIDAndAccountID(&c.ID, &c.AccountID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}
