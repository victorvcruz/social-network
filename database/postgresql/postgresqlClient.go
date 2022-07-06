package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgresqlClient struct {
	User     string
	Host     string
	Port     string
	Password string
	DbName   string
	Db       *sql.DB
}

func (p *PostgresqlClient) Conn() error {

	var DataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Password, p.DbName)

	db, err := sql.Open("postgres", DataSourceName)

	if err != nil {
		return err
	} else {
		log.Println("Connected to database " + p.DbName + "!")
	}

	p.Db = db

	return nil
}
