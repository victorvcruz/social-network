package post

import (
	"database/sql"
	"strings"
	"time"
)

type PostRepository interface {
	InsertPost(post *Post) error
	FindPostsByAccountID(accountID, page *string) ([]interface{}, error)
	UpdatePostDataByID(postID, accountID, content *string) error
	FindPostByID(id *string) (*PostResponse, error)
	ExistsPostByID(id *string) (*bool, error)
	RemovePostByID(postID, accountID *string) error
	ExistsPostByPostIDAndAccountID(postID, accountID *string) (*bool, error)
	FindPostByAccountFollowingByAccountID(accountID *string, page *string) ([]interface{}, error)
}

type PostRepositoryStruct struct {
	Db *sql.DB
}

func NewPostRepository(postgresDB *sql.DB) PostRepository {
	return &PostRepositoryStruct{postgresDB}
}

func (p *PostRepositoryStruct) InsertPost(post *Post) error {
	sqlStatement := `
		INSERT INTO post (id, account_id, content, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6)`

	row := p.Db.QueryRow(sqlStatement, post.ID, post.AccountID, post.Content, post.CreatedAt,
		post.UpdatedAt, post.Removed)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *PostRepositoryStruct) FindPostsByAccountID(accountID, page *string) ([]interface{}, error) {
	sqlStatement := `
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

	rows, err := p.Db.Query(sqlStatement, accountID, page)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var post PostResponse
	for rows.Next() {
		err = rows.Scan(
			&post.ID,
			&post.AccountID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Like,
			&post.Dislike,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = strings.Join(strings.Split(post.CreatedAt, "T00:00:00Z"), "")
		post.UpdatedAt = strings.Join(strings.Split(post.CreatedAt, "T00:00:00Z"), "")

		list = append(list, post)
	}

	return list, nil
}

func (p *PostRepositoryStruct) UpdatePostDataByID(postID, accountID, content *string) error {
	sqlStatement := `
		UPDATE post 
		SET content = $1, updated_at = $2
		WHERE id = $3
		AND account_id = $4  
		AND removed = false`

	updateTime := time.Now().UTC().Format("2006-01-02")

	row := p.Db.QueryRow(sqlStatement, content, updateTime, postID, accountID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil

}

func (p *PostRepositoryStruct) FindPostByID(id *string) (*PostResponse, error) {
	sqlStatement := `
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

	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var post PostResponse
	err = rows.Scan(
		&post.ID,
		&post.AccountID,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Like,
		&post.Dislike,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepositoryStruct) ExistsPostByID(id *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM post
		WHERE id = $1
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *PostRepositoryStruct) RemovePostByID(postID, accountID *string) error {
	sqlStatement := `
		UPDATE post 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	row := p.Db.QueryRow(sqlStatement, postID, accountID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *PostRepositoryStruct) ExistsPostByPostIDAndAccountID(postID, accountID *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM post
		WHERE id = $1
		AND account_id = $2
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, postID, accountID)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *PostRepositoryStruct) FindPostByAccountFollowingByAccountID(accountID, page *string) ([]interface{}, error) {
	sqlStatement := `
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

	rows, err := p.Db.Query(sqlStatement, accountID, page)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var post PostResponse
	for rows.Next() {
		err = rows.Scan(
			&post.ID,
			&post.AccountID,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Like,
			&post.Dislike,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = strings.Join(strings.Split(post.CreatedAt, "T00:00:00Z"), "")
		post.UpdatedAt = strings.Join(strings.Split(post.CreatedAt, "T00:00:00Z"), "")

		list = append(list, post)
	}

	return list, nil
}
