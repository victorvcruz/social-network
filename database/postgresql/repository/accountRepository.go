package repository

import (
	"database/sql"
	"social_network_project/entities"
)

type AccountRepository struct {
	Db *sql.DB
}

func (p *AccountRepository) InsertAccount(account *entities.Account) error {
	sqlStatement := `
		INSERT INTO account (id, username, name, description, email, password, created_at, updated_at, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := p.Db.Exec(sqlStatement, account.ID, account.Username, account.Name, account.Description,
		account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.Deleted)
	if err != nil {
		return err
	}

	return nil
}

func (p *AccountRepository) ExistsAccountByEmailAndPassword(email string, password string) bool {
	sqlStatement := `
		SELECT email, password 
		FROM account
		WHERE email = $1
		AND password = $2 `
	rows, _ := p.Db.Query(sqlStatement, email, password)

	return rows.Next()
}

func (p *AccountRepository) FindAccountIDbyEmail(email string) string {
	sqlStatement := `
		SELECT id
		FROM account
		WHERE email = $1`
	rows, _ := p.Db.Query(sqlStatement, email)
	rows.Next()
	var id string
	_ = rows.Scan(&id)

	return id
}
