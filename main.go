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

	router.HandleFunc("/api/servidores", pessoal.GetPessoal).Methods("GET")
	router.HandleFunc("/api/servidor/{matricula:[0-9]+}", pessoal.GetPessoalMat).Methods("GET") // URL parameter with Regex in URL

	log.Fatal(http.ListenAndServe(":8080", router))
}
