package main

import (
	"log"
	"net/http"
	"pessoalAPI-gorilamux/controllers"
	"pessoalAPI-gorilamux/db"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	db.Init()
	defer db.GetDB().Db.Close()
	pessoal := new(controllers.PessoalController)

	router.HandleFunc("/api/pessoal", pessoal.GetPessoal).Methods("GET")
	// router.HandleFunc("/resources", CreateResource).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
