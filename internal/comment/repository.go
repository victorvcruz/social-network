package comment

import (
	"database/sql"
	"strings"
	"time"
)

type CommentRepository interface {
	InsertComment(comment *Comment) error
	ExistsCommentByID(id *string) (*bool, error)
	FindCommentsByAccountID(accountID, page *string) ([]interface{}, error)
	FindCommentsByPostOrCommentID(postID, commentID, page *string) ([]interface{}, error)
	UpdateCommentDataByID(commentID, accountID, content *string) error
	FindCommentByID(id *string) (*Comment, error)
	RemoveCommentByID(commentID, accountID *string) error
	ExistsCommentByCommentIDAndAccountID(commentID, accountID *string) (*bool, error)
	FindAccountEmailOfPostByCommentID(commentID *string) ([]interface{}, error)
	FindAccountEmailOfPostAndCommentByCommentID(commentID *string) ([]interface{}, error)
}

type CommentRepositoryStruct struct {
	Db *sql.DB
}

func NewComentRepository(postgresDB *sql.DB) CommentRepository {
	return &CommentRepositoryStruct{postgresDB}
}

func (p *CommentRepositoryStruct) InsertComment(comment *Comment) error {
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

func (p *CommentRepositoryStruct) FindCommentsByAccountID(accountID, page *string) ([]interface{}, error) {

	stringQuery := `
		SELECT comment.id, comment.account_id, comment.post_id, comment.comment_id, comment.content, comment.created_at, comment.updated_at, 
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

	rows, err := p.Db.Query(stringQuery, accountID, page)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var comment Comment
	for rows.Next() {
		err = rows.Scan(
			&comment.ID,
			&comment.AccountID,
			&comment.PostID,
			&comment.CommentID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Like,
			&comment.Dislike,
		)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = strings.Join(strings.Split(comment.CreatedAt, "T00:00:00Z"), "")
		comment.UpdatedAt = strings.Join(strings.Split(comment.CreatedAt, "T00:00:00Z"), "")
		list = append(list, comment.ToResponse())
	}

	return list, nil
}

func (p *CommentRepositoryStruct) FindCommentsByPostOrCommentID(postID, commentID, page *string) ([]interface{}, error) {

	var str string
	var value *string
	if *postID != "" {
		str = "comment.post_id = $1 "
		value = postID
	} else {
		str = "comment.comment_id = $1 "
		value = commentID
	}

	stringQuery := `
		SELECT comment.id, comment.account_id, comment.post_id, comment.comment_id, comment.content, comment.created_at, comment.updated_at, 
	(
		SELECT count(1) FROM interaction i WHERE i.comment_id = comment.id AND i."type" = 'LIKE' 
	) AS like,
	(
		SELECT count(1) FROM interaction i WHERE i.comment_id = comment.id AND i."type" = 'DISLIKE' 
	) AS dislike
	FROM comment
	WHERE ` + str +
		`AND comment.removed = false
		Order By comment.created_at 
		OFFSET ($2 - 1) * 10
		FETCH NEXT 10 ROWS ONLY;`

	rows, err := p.Db.Query(stringQuery, value, page)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var comment Comment
	for rows.Next() {
		err = rows.Scan(
			&comment.ID,
			&comment.AccountID,
			&comment.PostID,
			&comment.CommentID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Like,
			&comment.Dislike,
		)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = strings.Join(strings.Split(comment.CreatedAt, "T00:00:00Z"), "")
		comment.UpdatedAt = strings.Join(strings.Split(comment.CreatedAt, "T00:00:00Z"), "")
		list = append(list, comment.ToResponse())
	}

	return list, nil
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

func (p *CommentRepositoryStruct) FindCommentByID(id *string) (*Comment, error) {
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

	var comment Comment
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

func (p *CommentRepositoryStruct) FindAccountEmailOfPostByCommentID(commentID *string) ([]interface{}, error) {
	sqlStatement := `
		SELECT account.email
		FROM comment
		INNER JOIN post ON post.id = comment.post_id 
		INNER JOIN account ON account.id = post.account_id
		WHERE comment.id = $1`

	rows, err := p.Db.Query(sqlStatement, commentID)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var emailOfPost string
	rows.Next()
	err = rows.Scan(
		&emailOfPost,
	)
	if err != nil {
		return nil, err
	}

	list = append(list, emailOfPost)

	return list, nil
}

func (p *CommentRepositoryStruct) FindAccountEmailOfPostAndCommentByCommentID(commentID *string) ([]interface{}, error) {
	sqlStatement := `
	SELECT account.email,	
	(
	SELECT account.email 
	FROM comment c
	INNER JOIN account ON account.id = comment.account_id
	WHERE c.id = comment.id
	) AS email2
	FROM comment
	INNER JOIN post ON post.id = comment.post_id 
	INNER JOIN account ON account.id = post.account_id
	WHERE comment.id = $1`

	rows, err := p.Db.Query(sqlStatement, commentID)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var emailOfPost string
	var emailOfComment string

	rows.Next()
	err = rows.Scan(
		&emailOfPost,
		&emailOfComment,
	)
	if err != nil {
		return nil, err
	}

	list = append(list, emailOfPost, emailOfComment)

	return list, nil
}
