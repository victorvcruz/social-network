package repository

import (
	"database/sql"
	"social_network_project/entities"
	"time"
)

type InteractionRepository interface {
	InsertInteraction(interaction *entities.Interaction) error
	ExistsInteractionByID(id *string) (*bool, error)
	UpdateInteractonDataByID(interactionID, accountID *string, typeValue *entities.InteractionType) error
	ExistsInteractionByInteractionIDAndAccountID(interactionID, accountID *string) (*bool, error)
	FindInteractionByID(id *string) (*entities.Interaction, error)
	RemoveInteractionByID(interactionID, accountID *string) error
	ExistsInteractionByPostIDAndAccountID(postID, accountID *string) (*bool, error)
	ExistsInteractionByCommentIDAndAccountID(commentID, accountID *string) (*bool, error)
	FindAccountEmailOfPostByInteractionID(interactionID *string) ([]interface{}, error)
	FindAccountEmailOfCommentByInteractionID(interactionID *string) ([]interface{}, error)
}

type InteractionRepositoryStruct struct {
	Db *sql.DB
}

func NewInteractionRepository(postgresDB *sql.DB) InteractionRepository {
	return &InteractionRepositoryStruct{postgresDB}
}

func (p *InteractionRepositoryStruct) InsertInteraction(interaction *entities.Interaction) error {
	sqlStatement := `
		INSERT INTO interaction (id, account_id, post_id, comment_id, type, created_at, updated_at, removed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	row := p.Db.QueryRow(sqlStatement, interaction.ID, interaction.AccountID, interaction.PostID, interaction.CommentID,
		interaction.Type.ToString(), interaction.CreatedAt, interaction.UpdatedAt, interaction.Removed)
	if row.Err() != nil {
		return row.Err()
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

func (p *InteractionRepositoryStruct) UpdateInteractonDataByID(interactionID, accountID *string, typeValue *entities.InteractionType) error {
	sqlStatement := `
		UPDATE interaction
		SET type = $1, updated_at = $2
		WHERE id = $3
		AND account_id = $4
		AND removed = false`

	updateTime := time.Now().UTC().Format("2006-01-02")

	row := p.Db.QueryRow(sqlStatement, typeValue.ToString(), updateTime, interactionID, accountID)
	if row.Err() != nil {
		return row.Err()
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
	var typeStr string
	err = rows.Scan(
		&interaction.ID,
		&interaction.AccountID,
		&interaction.PostID,
		&interaction.CommentID,
		&typeStr,
		&interaction.CreatedAt,
		&interaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	interaction.Type, _ = entities.ParseString(typeStr)

	return &interaction, nil
}

func (p *InteractionRepositoryStruct) RemoveInteractionByID(interactionID, accountID *string) error {
	sqlStatement := `
		UPDATE interaction 
		SET removed = true
		WHERE id = $1
		AND account_id = $2`

	row := p.Db.QueryRow(sqlStatement, interactionID, accountID)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (p *InteractionRepositoryStruct) ExistsInteractionByPostIDAndAccountID(postID, accountID *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM interaction
		WHERE post_id = $1
		AND account_id = $2
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, postID, accountID)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *InteractionRepositoryStruct) ExistsInteractionByCommentIDAndAccountID(commentID, accountID *string) (*bool, error) {
	sqlStatement := `
		SELECT id
		FROM interaction
		WHERE comment_id = $1
		AND account_id = $2
		AND removed = false`
	rows, err := p.Db.Query(sqlStatement, commentID, accountID)
	if err != nil {
		return nil, err
	}

	next := rows.Next()
	return &next, nil
}

func (p *InteractionRepositoryStruct) FindAccountEmailOfPostByInteractionID(interactionID *string) ([]interface{}, error) {
	sqlStatement := `
	SELECT account.email
	FROM interaction
	INNER JOIN post ON post.id = interaction.post_id 
	INNER JOIN account ON account.id = post.account_id
	WHERE interaction.id = $1`

	rows, err := p.Db.Query(sqlStatement, interactionID)
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

func (p *InteractionRepositoryStruct) FindAccountEmailOfCommentByInteractionID(interactionID *string) ([]interface{}, error) {
	sqlStatement := `
	SELECT account.email
	FROM interaction
	INNER JOIN comment ON comment.id = interaction.comment_id  
	INNER JOIN account ON account.id = comment.account_id
	WHERE interaction.id = $1`

	rows, err := p.Db.Query(sqlStatement, interactionID)
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
