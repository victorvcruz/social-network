package repository

import (
	"database/sql"
	"social_network_project/database/postgresql"
)

type PostRepository interface {
}

type PostRepositoryStruct struct {
	Db *sql.DB
}

func NewPostRepository() PostRepository {
	return &PostRepositoryStruct{postgresql.Db}
}
