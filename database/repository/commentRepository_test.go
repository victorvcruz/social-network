package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"social_network_project/entities/response"
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

	ids := map[string]interface{}{
		"post_id":    "7d077271-fc44-4792-9ba9-ff15154f6cca",
		"comment_id": "05f508c6-cc9b-4942-a62f-0aa08d01ed45",
	}

	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := dinamicQueryFindCommentsByAccountID(ids)

	rows := sqlmock.NewRows([]string{"id", "account_id", "post_id", "comment_id", "content", "created_at", "updated_at"}).
		AddRow(c.ID, c.AccountID, c.PostID, c.CommentID, c.Content, p.CreatedAt, p.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(c.ID).WillReturnRows(rows)

	account, err := repo.FindCommentsByAccountID(&c.AccountID, &c.PostID, &c.CommentID.String)
	assert.Empty(t, account)
	assert.Error(t, err)
}

func TestNewComentRepository_dinamicQueryFindCommentsByAccountID(t *testing.T) {
	ids1 := map[string]interface{}{
		"post_id":    "7d077271-fc44-4792-9ba9-ff15154f6cca",
		"comment_id": "05f508c6-cc9b-4942-a62f-0aa08d01ed45",
	}

	stringQueryExpected1 := `SELECT id, account_id, post_id, comment_id, content, created_at, updated_at FROM comment WHERE account_id = $1 AND removed = false AND "comment_id" = '05f508c6-cc9b-4942-a62f-0aa08d01ed45' AND "post_id" = '7d077271-fc44-4792-9ba9-ff15154f6cca'`
	stringQuery1 := dinamicQueryFindCommentsByAccountID(ids1)

	assert.Equal(t, len(stringQueryExpected1), len(stringQuery1))

	ids2 := map[string]interface{}{
		"post_id":    "",
		"comment_id": "",
	}

	stringQueryExpected2 := `SELECT id, account_id, post_id, comment_id, content, created_at, updated_at FROM comment WHERE account_id = $1 AND removed = false`
	stringQuery2 := dinamicQueryFindCommentsByAccountID(ids2)

	assert.Equal(t, len(stringQueryExpected2), len(stringQuery2))
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

func TestCommentRepositoryStruct_CountInteractionsForComment(t *testing.T) {
	db, mock := NewMock()
	repo := CommentRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT count(type) 
		FROM interaction
		WHERE comment_id = $1
		AND type = $2`

	rows := sqlmock.NewRows([]string{"comment_id", "type"}).
		AddRow(i.CommentID, i.Type)

	mock.ExpectQuery(query).WithArgs(i.CommentID).WillReturnRows(rows)

	id, err := repo.CountInteractionsForComment(&i.CommentID.String, response.INTERACTION_TYPE_LIKED.EnumIndex())
	assert.Empty(t, id)
	assert.Error(t, err)
}
