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
	Change controllers.Change
	Delete controllers.Delete
}

func (a *Api) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/accounts", a.Create.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/auth", a.Create.CreateToken).Methods("POST")
	router.HandleFunc("/accounts", a.Read.GetAccount).Methods("GET")
	router.HandleFunc("/accounts", a.Change.ChangeAccount).Methods("PUT")
	router.HandleFunc("/accounts", a.Delete.DeleteAccount).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
