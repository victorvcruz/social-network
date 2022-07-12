package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"testing"
	"time"
)

var i = &entities.Interaction{
	ID:        uuid.New().String(),
	AccountID: "f981d822-7efb-4e66-aa84-99f517820ca3",
	PostID:    entities.NewNullString("0d0bb472-225c-4c8a-9935-a21045c80d87"),
	CommentID: entities.NewNullString("8b607c43-0190-4c8c-9746-4b527d1d2c55"),
	Type:      0,
	CreatedAt: time.Now().UTC().Format("2006-01-02"),
	UpdatedAt: time.Now().UTC().Format("2006-01-02"),
	Removed:   false,
}

func TestInteractionRepositoryStruct_InsertInteraction(t *testing.T) {
	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		INSERT INTO interaction (id, account_id, post_id, comment_id, type, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(i.ID, i.AccountID, i.PostID, i.CommentID, i.Type, i.CreatedAt,
		i.UpdatedAt, i.Removed).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.InsertInteraction(i)
	assert.Error(t, err)
}

func TestInteractionRepositoryStruct_ExistsInteractionByID(t *testing.T) {
	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM interaction
		WHERE id = $1
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(i.ID)

	mock.ExpectQuery(query).WithArgs(i.ID).WillReturnRows(rows)

	exist, err := repo.ExistsInteractionByID(&i.ID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestCommentRepositoryStruct_UpdateInteractonDataByID(t *testing.T) {

	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

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

	prep.ExpectExec().WithArgs(i.Type, i.UpdatedAt, i.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateInteractonDataByID(&i.ID, &i.AccountID, &i.Type)
	assert.Error(t, err)
}

func TestInteractionRepositoryStruct_ExistsInteractionByInteractionIDAndAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM interaction
		WHERE id = $1
		AND account_id = $2
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(i.ID)

	mock.ExpectQuery(query).WithArgs(i.ID).WillReturnRows(rows)

	exist, err := repo.ExistsInteractionByInteractionIDAndAccountID(&i.ID, &i.AccountID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestInteractionRepositoryStruct_FindInteractionsByID(t *testing.T) {
	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id, account_id, post_id, comment_id, type, created_at, updated_at
		FROM interaction
		WHERE id = $1
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id", "account_id", "type", "created_at", "updated_at"}).
		AddRow(i.ID, i.AccountID, i.Type, i.CreatedAt, i.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(i.ID).WillReturnRows(rows)

	account, err := repo.FindInteractionByID(&i.ID)
	assert.Empty(t, account)
	assert.Error(t, err)
}

func TestInteractionRepositoryStruct_RemoveInteractionByID(t *testing.T) {

	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		UPDATE interaction 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(i.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.RemoveInteractionByID(&i.ID, &i.AccountID)
	assert.Error(t, err)
}

func TestInteractionRepositoryStruct_ExistsInteractionByPostIDAndAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM interaction
		WHERE post_id = $1
		AND account_id = $2
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(i.ID)

	mock.ExpectQuery(query).WithArgs(i.PostID).WillReturnRows(rows)

	exist, err := repo.ExistsInteractionByInteractionIDAndAccountID(&i.PostID.String, &i.AccountID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestInteractionRepositoryStruct_ExistsInteractionByCommentIDAndAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := InteractionRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM interaction
		WHERE comment_id = $1
		AND account_id = $2
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(i.ID)

	mock.ExpectQuery(query).WithArgs(i.CommentID).WillReturnRows(rows)

	exist, err := repo.ExistsInteractionByCommentIDAndAccountID(&i.CommentID.String, &i.AccountID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}
