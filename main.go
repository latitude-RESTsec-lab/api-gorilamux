package main

import (
	"log"
	"net/http"

	"github.com/latitude-RESTsec-lab/api-gorilamux/db"

	"github.com/latitude-RESTsec-lab/api-gorilamux/controllers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	db.Init()
	defer db.GetDB().Db.Close()
	Servidor := new(controllers.ServidorController)

	router.HandleFunc("/api/servidores", Servidor.GetServidor).Methods("GET")
	router.HandleFunc("/api/servidor/", Servidor.PostServidor).Methods("POST")
	router.HandleFunc("/api/servidor/{matricula:[0-9]+}", Servidor.GetServidorMat).Methods("GET") // URL parameter with Regex in URL

	log.Fatal(http.ListenAndServe(":8080", router))
}
