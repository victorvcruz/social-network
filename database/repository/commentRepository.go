package repository

import (
	"database/sql"
	"social_network_project/database/postgresql"
	"social_network_project/entities"
)

type CommentRepository interface {
	InsertComment(comment *entities.Comment) error
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

	_, err := p.Db.Exec(sqlStatement, comment.ID, comment.AccountID, comment.PostID, comment.CommentID,
		comment.Content, comment.CreatedAt, comment.UpdatedAt, comment.Removed)
	if err != nil {
		return err
	}

	return nil
}
