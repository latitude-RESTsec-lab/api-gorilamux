package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/restsec/api-gorilamux/db"

	"github.com/restsec/api-gorilamux/config"
	"github.com/restsec/api-gorilamux/controllers"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting GORILLA MUX API")
	err := config.ReadConfig()
	if err != nil {
		log.Print(err.Error())
		return
	}

	log.SetOutput(config.LogFile)

	db.Init()
	defer db.GetDB().Db.Close()

	Servidor := new(controllers.ServidorController)

	httpsRouter := mux.NewRouter()

	httpsRouter.HandleFunc("/api/servidores", Servidor.GetServidor).Methods("GET")
	httpsRouter.HandleFunc("/api/servidor/", Servidor.PostServidor).Methods("POST")
	httpsRouter.HandleFunc("/api/servidor/{matricula:[0-9]+}", Servidor.GetServidorMat).Methods("GET") // URL parameter with Regex in URL
	httpsRouter.HandleFunc("/api/calculo/", Servidor.Calculate).Methods("POST")
	err = http.ListenAndServeTLS(":"+config.ConfigParams.HttpsPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation, httpsRouter)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}

}
