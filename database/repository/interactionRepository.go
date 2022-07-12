package repository

import (
	"database/sql"
	"social_network_project/database/postgresql"
	"social_network_project/entities"
	"social_network_project/entities/response"
	"time"
)

type InteractionRepository interface {
	InsertInteraction(interaction *entities.Interaction) error
	ExistsInteractionByID(id *string) (*bool, error)
	UpdateInteractonDataByID(interactionID, accountID *string, typeValue *response.InteractionType) error
	ExistsInteractionByInteractionIDAndAccountID(interactionID, accountID *string) (*bool, error)
	FindInteractionByID(id *string) (*entities.Interaction, error)
	RemoveInteractionByID(interactionID, accountID *string) error
}

type InteractionRepositoryStruct struct {
	Db *sql.DB
}

func NewInteractionRepository() InteractionRepository {
	return &InteractionRepositoryStruct{postgresql.Db}
}

func (p *InteractionRepositoryStruct) InsertInteraction(interaction *entities.Interaction) error {
	sqlStatement := `
		INSERT INTO interaction (id, account_id, post_id, comment_id, type, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := p.Db.Exec(sqlStatement, interaction.ID, interaction.AccountID, interaction.PostID, interaction.CommentID,
		interaction.Type, interaction.CreatedAt, interaction.UpdatedAt, interaction.Removed)
	if err != nil {
		return err
	}

	return nil
}

func (p *InteractionRepositoryStruct) ExistsInteractionByID(id *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM interaction
		WHERE id = $1
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *InteractionRepositoryStruct) UpdateInteractonDataByID(interactionID, accountID *string, typeValue *response.InteractionType) error {
	sqlStatement := `
		UPDATE interaction
		SET type = $1, updated_at = $2
		WHERE id = $3
		AND account_id = $4
		AND removed = false`

	updateTime := time.Now().UTC().Format("2006-01-02")

	_, err := p.Db.Exec(sqlStatement, typeValue, updateTime, interactionID, accountID)
	if err != nil {
		return err
	}

	return nil
}

func (p *InteractionRepositoryStruct) ExistsInteractionByInteractionIDAndAccountID(interactionID, accountID *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM interaction
		WHERE id = $1
		AND account_id = $2
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, interactionID, accountID)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *InteractionRepositoryStruct) FindInteractionByID(id *string) (*entities.Interaction, error) {
	sqlStatement := `
		SELECT id, account_id, post_id, comment_id, type, created_at, updated_at
		FROM interaction
		WHERE id = $1
		AND removed = false`

	rows, err := p.Db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	rows.Next()

	var interaction entities.Interaction
	err = rows.Scan(
		&interaction.ID,
		&interaction.AccountID,
		&interaction.PostID,
		&interaction.CommentID,
		&interaction.Type,
		&interaction.CreatedAt,
		&interaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &interaction, nil
}

func (p *InteractionRepositoryStruct) RemoveInteractionByID(interactionID, accountID *string) error {
	sqlStatement := `
		UPDATE interaction 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	_, err := p.Db.Exec(sqlStatement, interactionID, accountID)
	if err != nil {
		return err
	}

	return nil
}
