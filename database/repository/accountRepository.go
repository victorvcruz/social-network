package repository

import (
	"database/sql"
	"fmt"
	"social_network_project/database/postgresql"
	"social_network_project/entities"
	"strings"
)

type AccountRepository interface {
	InsertAccount(account *entities.Account) error
	FindAccountPasswordByEmail(email string) (*string, error)
	FindAccountIDbyEmail(email string) (*string, error)
	FindAccountByID(id *string) (*entities.Account, error)
	ChangeAccountDataByID(id *string, mapBody map[string]interface{}) error
	DeleteAccountByID(id *string) error
	ExistsAccountByID(id *string) (*bool, error)
	ExistsAccountByUsername(username *string) (*bool, error)
	ExistsAccountByEmail(email *string) (*bool, error)
	InsertAccountFollow(accountID, accountFollow *string) error
	FindAccountFollowingByAccountID(accountID *string) ([]interface{}, error)
	FindAccountFollowersByAccountID(accountID *string) ([]interface{}, error)
	ExistsFollowByAccountIDAndAccountFollowedID(accountID, accountToFollow *string) (*bool, error)
	DeleteAccountFollow(accountID, accountFollow *string) error
}

type AccountRepositoryStruct struct {
	Db *sql.DB
}

func NewAccountRepository() AccountRepository {
	return &AccountRepositoryStruct{postgresql.Db}
}

func (p *AccountRepositoryStruct) InsertAccount(account *entities.Account) error {
	sqlStatement := `
		INSERT INTO account (id, username, name, description, email, password, created_at, updated_at, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	row := p.Db.QueryRow(sqlStatement, account.ID, account.Username, account.Name, account.Description,
		account.Email, account.Password, account.CreatedAt, account.UpdatedAt, account.Deleted)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *AccountRepositoryStruct) FindAccountPasswordByEmail(email string) (*string, error) {
	sqlStatement := `
		SELECT password 
		FROM account
		WHERE email = $1
		AND deleted = false`
	rows, err := p.Db.Query(sqlStatement, email)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var password *string
	_ = rows.Scan(&password)

	return password, nil
}

func (p *AccountRepositoryStruct) FindAccountIDbyEmail(email string) (*string, error) {
	sqlStatement := `
		SELECT id
		FROM account
		WHERE email = $1
		AND deleted = false`
	rows, err := p.Db.Query(sqlStatement, email)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var id *string
	err = rows.Scan(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (p *AccountRepositoryStruct) FindAccountByID(id *string) (*entities.Account, error) {
	sqlStatement := `
		SELECT id, username, name, description, email, password, created_at, updated_at, deleted
		FROM account
		WHERE id = $1
		AND deleted = false`
	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var account entities.Account
	err = rows.Scan(
		&account.ID,
		&account.Username,
		&account.Name,
		&account.Description,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.Deleted,
	)
	if err != nil {
		return nil, err
	}

	account.CreatedAt = strings.Join(strings.Split(account.CreatedAt, "T00:00:00Z"), "")
	account.UpdatedAt = strings.Join(strings.Split(account.CreatedAt, "T00:00:00Z"), "")

	return &account, nil
}

func (p *AccountRepositoryStruct) ChangeAccountDataByID(id *string, mapBody map[string]interface{}) error {
	sqlStatement := dinamicQueryChangeAccountDataByID(mapBody)

	row := p.Db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func dinamicQueryChangeAccountDataByID(mapBody map[string]interface{}) string {

	var values []interface{}
	var where []string

	for key, value := range mapBody {
		values = append(values, value)
		where = append(where, fmt.Sprintf(`"%s" = '%s'`, key, value))
	}
	stringQuery := "UPDATE account SET " + strings.Join(where, ", ") + " WHERE id = $1 AND deleted = false"

	return stringQuery
}

func (p *AccountRepositoryStruct) DeleteAccountByID(id *string) error {
	sqlStatement := `
		UPDATE account 
		SET deleted = true 
		WHERE id = $1
		AND deleted = false`

	row := p.Db.QueryRow(sqlStatement, id)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *AccountRepositoryStruct) ExistsAccountByID(id *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM account
		WHERE id = $1
		AND deleted = false`
	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *AccountRepositoryStruct) ExistsAccountByUsername(username *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM account
		WHERE username = $1
		AND deleted = false`
	rows, err := p.Db.Query(sqlStatement, username)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *AccountRepositoryStruct) ExistsAccountByEmail(email *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM account
		WHERE email = $1
		AND deleted = false`
	rows, err := p.Db.Query(sqlStatement, email)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *AccountRepositoryStruct) InsertAccountFollow(accountID, accountFollow *string) error {
	sqlStatement := `
		INSERT INTO account_follow (account_id, account_id_followed, unfollowed)
		VALUES ($1, $2, false)`

	row := p.Db.QueryRow(sqlStatement, accountID, accountFollow)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *AccountRepositoryStruct) FindAccountFollowingByAccountID(accountID *string) ([]interface{}, error) {
	sqlStatement := `
		SELECT account_id_followed
		FROM account_follow
		WHERE account_id = $1
		AND unfollowed = false`

	rows, err := p.Db.Query(sqlStatement, accountID)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var id *string

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		account, err := p.FindAccountByID(id)
		if err != nil {
			return nil, err
		}

		list = append(list, account.ToResponse())
	}

	return list, nil
}

func (p *AccountRepositoryStruct) FindAccountFollowersByAccountID(accountID *string) ([]interface{}, error) {
	sqlStatement := `
		SELECT account_id
		FROM account_follow
		WHERE account_id_followed = $1
		AND unfollowed = false`

	rows, err := p.Db.Query(sqlStatement, accountID)
	if err != nil {
		return nil, err
	}

	list := []interface{}{}
	var id *string

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		account, err := p.FindAccountByID(id)
		if err != nil {
			return nil, err
		}

		list = append(list, account.ToResponse())
	}

	return list, nil
}

func (p *AccountRepositoryStruct) DeleteAccountFollow(accountID, accountFollow *string) error {
	sqlStatement := `
		UPDATE account_follow 
		SET unfollowed = true 
		WHERE account_id= $1
		AND account_id_followed = $2
		AND unfollowed = false`

	row := p.Db.QueryRow(sqlStatement, accountID, accountFollow)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *AccountRepositoryStruct) ExistsFollowByAccountIDAndAccountFollowedID(accountID, accountToFollow *string) (*bool, error) {
	sqlStatement := `
		SELECT account_id
		FROM account_follow
		WHERE account_id = $1
		AND account_id_followed = $2
		AND unfollowed = false`
	rows, err := p.Db.Query(sqlStatement, accountID, accountToFollow)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}
