package main

import (
	"log"
	"social_network_project/api"
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
	read := controllers.Read{
		AccountRepository: accountRepository,
	}

	change := controllers.Change{
		AccountRepository: accountRepository,
	}

	api := api.Api{
		Create: create,
		Read:   read,
		Change: change,
	}

	api.Run()

}
