package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"social_network_project/utils"
)

func ConnectDatabase() (*sql.DB, error) {

	dbUser := utils.GetStringEnvOrElse("POSTGRESQL_USER", "postgres")
	dbPwd := utils.GetStringEnvOrElse("POSTGRESQL_PASSWORD", "admin")
	DBName := utils.GetStringEnvOrElse("DB_NAME", "postgres")
	dbHost := utils.GetStringEnvOrElse("DB_HOST", "localhost")
	dbPort := utils.GetStringEnvOrElse("DB_PORT", "5432")
	postgresURI := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPwd, DBName)

	db, err := sql.Open("postgres", postgresURI)

	if err != nil {
		return nil, err
	} else {
		log.Println("Connected to database " + DBName + "!")
	}

	return db, nil
}
