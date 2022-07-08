package main

import (
	"github.com/go-playground/validator/v10"
	"log"
	"social_network_project/api"
	"social_network_project/api/handler"
	"social_network_project/controllers"
	"social_network_project/database/postgresql"
	"social_network_project/database/repository"
)

func main() {
	postgresqlClient := postgresql.PostgresqlClient{
		User:     "postgres",
		Host:     "localhost",
		Port:     "5432",
		Password: "postgres",
		DbName:   "postgres",
	}

	if err := postgresqlClient.Conn(); err != nil {
		log.Fatal(err)
	}

	accountRepository := repository.AccountRepository{
		Db: postgresqlClient.Db,
	}

	create := controllers.Create{
		AccountRepository: accountRepository,
	}

	AccountsRequest := request.AccountsAPI{
		AccountRepository: accountRepository,
		Validate:          validator.New(),
		Create:            create,
	}

	api := api.Api{
		AccountsAPI: AccountsRequest,
	}

	api.Run()

}
