package main

import (
	"log"
	"social_network_project/api"
	"social_network_project/controllers"
	"social_network_project/database/postgresql"
	"social_network_project/database/postgresql/repository"
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
	view := controllers.View{}

	api := api.Api{
		Create: create,
		View:   view,
	}

	api.Run()

}
