package repository

import (
	"database/sql"
	"fmt"
	"social_network_project/database/postgresql"
	"social_network_project/entities"
	"social_network_project/entities/response"
	"strings"
	"time"
)

type CommentRepository interface {
	InsertComment(comment *entities.Comment) error
	ExistsCommentByID(id *string) (*bool, error)
	FindCommentsByAccountID(accountID, postID, commentID *string) ([]interface{}, error)
	UpdateCommentDataByID(commentID, accountID, content *string) error
	FindCommentByID(id *string) (*entities.Comment, error)
	RemoveCommentByID(commentID, accountID *string) error
	ExistsCommentByCommentIDAndAccountID(commentID, accountID *string) (*bool, error)
	CountInteractionsForComment(commentId *string, typeValue int) (*int, error)
}

type CommentRepositoryStruct struct {
	Db *sql.DB
}

func NewComentRepository() CommentRepository {
	return &CommentRepositoryStruct{postgresql.Db}
}

func (p *CommentRepositoryStruct) InsertComment(comment *entities.Comment) error {
	sqlStatement := `
		INSERT INTO comment (id, account_id, post_id, comment_id, content, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	row := p.Db.QueryRow(sqlStatement, comment.ID, comment.AccountID, comment.PostID, comment.CommentID,
		comment.Content, comment.CreatedAt, comment.UpdatedAt, comment.Removed)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *CommentRepositoryStruct) ExistsCommentByID(id *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM comment
		WHERE id = $1
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *CommentRepositoryStruct) FindCommentsByAccountID(accountID, postID, commentID *string) ([]interface{}, error) {

	ids := map[string]interface{}{
		"post_id":    *postID,
		"comment_id": *commentID,
	}

	str := dinamicQueryFindCommentsByAccountID(ids)

	rows, err := p.Db.Query(str, accountID)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var comment entities.Comment
	for rows.Next() {
		err = rows.Scan(
			&comment.ID,
			&comment.AccountID,
			&comment.PostID,
			&comment.CommentID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Removed,
		)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = strings.Join(strings.Split(comment.CreatedAt, "T00:00:00Z"), "")
		comment.UpdatedAt = strings.Join(strings.Split(comment.CreatedAt, "T00:00:00Z"), "")

		commentResponse := comment.ToResponse()

		like, err := p.CountInteractionsForComment(&commentResponse.ID, response.INTERACTION_TYPE_LIKED.EnumIndex())
		if err != nil {
			return nil, err
		}

		dislike, err := p.CountInteractionsForComment(&commentResponse.ID, response.INTERACTION_TYPE_DISLIKED.EnumIndex())
		if err != nil {
			return nil, err
		}

		commentResponse.Like = *like
		commentResponse.Dislike = *dislike
		list = append(list, commentResponse)
	}

	return list, nil
}

func dinamicQueryFindCommentsByAccountID(mapBody map[string]interface{}) string {

	var values []interface{}
	var where []string

	for key, value := range mapBody {
		values = append(values, value)
		if value != "" {
			where = append(where, fmt.Sprintf(`"%s" = '%s'`, key, value))
		}
	}
	str := ""
	if strings.Join(where, " AND ") != "" {
		str = " AND " + strings.Join(where, " AND ")
	}
	stringQuery := "SELECT id, account_id, post_id, comment_id, content, created_at, updated_at, removed FROM comment WHERE account_id = $1 AND removed = false" + str

	return stringQuery
}

func (p *CommentRepositoryStruct) UpdateCommentDataByID(commentID, accountID, content *string) error {
	sqlStatement := `
		UPDATE comment
		SET content = $1, updated_at = $2
		WHERE id = $3
		AND account_id = $4
		AND removed = false`

	updateTime := time.Now().UTC().Format("2006-01-02")

	row := p.Db.QueryRow(sqlStatement, content, updateTime, commentID, accountID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *CommentRepositoryStruct) FindCommentByID(id *string) (*entities.Comment, error) {
	sqlStatement := `
		SELECT id, account_id, post_id, comment_id, content, created_at, updated_at
		FROM comment
		WHERE id = $1
		AND removed = false`

	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	rows.Next()

	var comment entities.Comment
	err = rows.Scan(
		&comment.ID,
		&comment.AccountID,
		&comment.PostID,
		&comment.CommentID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (p *CommentRepositoryStruct) RemoveCommentByID(commentID, accountID *string) error {
	sqlStatement := `
		UPDATE comment 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	row := p.Db.QueryRow(sqlStatement, commentID, accountID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *CommentRepositoryStruct) ExistsCommentByCommentIDAndAccountID(commentID, accountID *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM comment
		WHERE id = $1
		AND account_id = $2
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, commentID, accountID)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *CommentRepositoryStruct) CountInteractionsForComment(commentId *string, typeValue int) (*int, error) {
	sqlStatement := `
		SELECT count(type) 
		FROM interaction
		WHERE comment_id = $1
		AND type = $2`

	rows, err := p.Db.Query(sqlStatement, commentId, typeValue)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var count *int
	err = rows.Scan(&count)
	if err != nil {
		return nil, err
	}

	return count, nil
}
