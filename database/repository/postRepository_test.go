package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"social_network_project/entities"
	"testing"
	"time"
)

var p = &entities.Post{
	ID:        uuid.New().String(),
	AccountID: "f981d822-7efb-4e66-aa84-99f517820ca3",
	Content:   "Marcelo Sabido",
	CreatedAt: time.Now().UTC().Format("2006-01-02"),
	UpdatedAt: time.Now().UTC().Format("2006-01-02"),
	Removed:   false,
}

var pBody = map[string]interface{}{
	"content": "eu sou baianao",
}

func TestPostRepositoryStruct_InsertPost(t *testing.T) {
	db, mock := NewMock()
	repo := PostRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		INSERT INTO post (id, account_id, content, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6)`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(p.ID, p.AccountID, p.Content, p.CreatedAt,
		p.UpdatedAt, p.Removed).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.InsertPost(p)
	assert.Error(t, err)
}

func TestPostRepositoryStruct_FindPostsByAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := PostRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	page := "1"
	query := `
	SELECT post.id, post.account_id, post.content, post.created_at, post.updated_at, 
	(
		SELECT count(1) FROM interaction i WHERE i.post_id = post.id AND i."type" = 'LIKE' 
	) AS like,
	(
		SELECT count(1) FROM interaction i WHERE i.post_id = post.id AND i."type" = 'DISLIKE' 
	) AS dislike
	FROM post
	WHERE post.account_id = $1
	AND post.removed = false
	Order By post.created_at 
	OFFSET ($2 - 1) * 10
	FETCH NEXT 10 ROWS ONLY;`

	rows := sqlmock.NewRows([]string{"id", "account_id", "content", "created_at", "updated_at"}).
		AddRow(p.ID, p.AccountID, p.Content, p.CreatedAt, p.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(p.ID).WillReturnRows(rows)

	account, err := repo.FindPostsByAccountID(&p.ID, &page)
	assert.Empty(t, account)
	assert.Error(t, err)
}

func TestPostRepositoryStruct_UpdatePostDataByID(t *testing.T) {

	db, mock := NewMock()
	repo := PostRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		UPDATE post 
		SET content = $1, updated_at = $2
		WHERE id = $3
		AND account_id = $4  
		AND removed = false`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(p.Content, p.UpdatedAt, p.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdatePostDataByID(&p.ID, &p.AccountID, &p.Content)
	assert.Error(t, err)
}

func TestPostRepositoryStruct_FindPostByID(t *testing.T) {
	db, mock := NewMock()
	repo := PostRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
	SELECT post.id, post.account_id, post.content, post.created_at, post.updated_at, 
	(
		SELECT count(1) FROM interaction i WHERE i.post_id = post.id AND i."type" = 'LIKE' 
	) AS like,
	(
		SELECT count(1) FROM interaction i WHERE i.post_id = post.id AND i."type" = 'DISLIKE' 
	) AS dislike
	FROM post
	WHERE post.id = $1
	AND post.removed = false
	GROUP BY post.id;`

	rows := sqlmock.NewRows([]string{"id", "account_id", "content", "created_at", "updated_at"}).
		AddRow(p.ID, p.AccountID, p.Content, p.CreatedAt, p.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(p.ID).WillReturnRows(rows)

	account, err := repo.FindPostByID(&p.ID)
	assert.Empty(t, account)
	assert.Error(t, err)
}

func TestPostRepositoryStruct_ExistsPostByID(t *testing.T) {
	db, mock := NewMock()
	repo := PostRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM post
		WHERE id = $1
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(p.ID)

	mock.ExpectQuery(query).WithArgs(p.ID).WillReturnRows(rows)

	exist, err := repo.ExistsPostByID(&p.ID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestPostRepositoryStruct_RemovePostByID(t *testing.T) {

	db, mock := NewMock()
	repo := PostRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		UPDATE post 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	prep := mock.ExpectPrepare(query)

	prep.ExpectExec().WithArgs(p.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.RemovePostByID(&p.ID, &p.AccountID)
	assert.Error(t, err)
}

func TestPostRepositoryStruct_ExistsPostByPostIDAndAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := PostRepositoryStruct{db}

	defer func() {
		db.Close()
	}()

	query := `
		SELECT id
		FROM post
		WHERE id = $1
		AND account_id = $2
		AND removed = false`

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(p.ID)

	mock.ExpectQuery(query).WithArgs(p.ID).WillReturnRows(rows)

	exist, err := repo.ExistsPostByPostIDAndAccountID(&p.ID, &p.AccountID)
	assert.Empty(t, exist)
	assert.Error(t, err)
}

func TestPostRepositoryStruct_FindPostByAccountFollowingByAccountID(t *testing.T) {
	db, mock := NewMock()
	repo := PostRepositoryStruct{db}
	page := "1"

	defer func() {
		db.Close()
	}()

	query := `
	SELECT post.id, post.account_id, post.content, post.created_at, post.updated_at, 
	(
		SELECT count(1) FROM interaction i WHERE i.post_id = post.id AND i."type" = 'LIKE' 
	) AS like,
	(
		SELECT count(1) FROM interaction i WHERE i.post_id = post.id AND i."type" = 'DISLIKE' 
	) AS dislike
	FROM account_follow
	INNER JOIN post ON account_follow.account_id_followed = post.account_id 
	WHERE account_follow.account_id = $1
	AND post.removed = false
	AND account_follow.unfollowed = false
	Order By post.created_at 
	OFFSET ($2 - 1) * 10
	FETCH NEXT 10 ROWS ONLY;`

	rows := sqlmock.NewRows([]string{"id", "account_id", "content", "created_at", "updated_at"}).
		AddRow(p.ID, p.AccountID, p.Content, p.CreatedAt, p.UpdatedAt)

	mock.ExpectQuery(query).WithArgs(p.ID).WillReturnRows(rows)

	account, err := repo.FindPostByAccountFollowingByAccountID(&p.AccountID, &page)
	assert.Empty(t, account)
	assert.Error(t, err)
}
