package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/latitude-RESTsec-lab/api-gorilamux/db"

	"github.com/latitude-RESTsec-lab/api-gorilamux/controllers"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting GORILLA MUX API")
	file, fileErr := os.Create("server.log")
	if fileErr != nil {
		fmt.Println(fileErr)
		file = os.Stdout
	}
	log.SetOutput(file)

	db.Init()
	defer db.GetDB().Db.Close()

	Servidor := new(controllers.ServidorController)

	httpsRouter := mux.NewRouter()

	httpsRouter.HandleFunc("/api/servidores", Servidor.GetServidor).Methods("GET")
	httpsRouter.HandleFunc("/api/servidor/", Servidor.PostServidor).Methods("POST")
	httpsRouter.HandleFunc("/api/servidor/{matricula:[0-9]+}", Servidor.GetServidorMat).Methods("GET") // URL parameter with Regex in URL
	err := http.ListenAndServeTLS(":443", "./devssl/server.pem", "./devssl/server.key", httpsRouter)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}

}
