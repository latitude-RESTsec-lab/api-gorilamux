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

	httpRouter := mux.NewRouter()
	httpsRouter.HandleFunc("/api/servidores", func(w http.ResponseWriter, r *http.Request) {
		target := "https://" + r.Host + r.URL.Path
		if len(r.URL.RawQuery) > 0 {
			target += "?" + r.URL.RawQuery
		}
		log.Printf("redirect to: %s", target)
		http.Redirect(w, r, target,
			http.StatusTemporaryRedirect)
	}).Methods("GET")
	httpsRouter.HandleFunc("/api/servidor/{matricula:[0-9]+}", Servidor.GetServidorMat).Methods("GET") // URL parameter with Regex in URL
	// http.DefaultTransport.RoundTrip("http2")
	go http.ListenAndServe(":80", httpRouter)
	log.Fatal(http.ListenAndServeTLS(":443", "./devssl/server.pem", "./devssl/server.key", httpsRouter))

}
