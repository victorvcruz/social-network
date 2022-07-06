package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"social_network_project/controllers"
)

type Api struct {
	Create controllers.Create
	Read   controllers.Read
}

func (a *Api) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", a.Create.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/auth", a.Create.CreateToken).Methods("POST")
	router.HandleFunc("/accounts", a.Read.GetAccount).Methods("GET")
	router.HandleFunc("/accounts", controllers.ChangeAccount).Methods("PUT")
	router.HandleFunc("/accounts", controllers.DeleteAccount).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
