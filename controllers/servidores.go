package controllers

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/latitude-RESTsec-lab/api-gorilamux/db"
)

//ErrorBody structure is used to improve error reporting in a JSON response body
type ErrorBody struct {
	Reason string `json:"reason"`
}

type Result struct {
	Result float64 `json:"Result"`
}

//Servidor structure is used to store data used by this API
type Servidor struct {
	ID                int    `db:"id, primarykey, autoincrement" json:"id"`
	Siape             int    `db:"siape" json:"siape"`
	Idpessoa          int    `db:"id_pessoa" json:"id_pessoa"`
	Nome              string `db:"nome" json:"nome"`
	Matriculainterna  int    `db:"matricula_interna" json:"matricula_interna"`
	Nomeidentificacao string `db:"nome_identificacao" json:"nome_identificacao"`
	Datanascimento    string `db:"data_nascimento" json:"data_nascimento"`
	Sexo              string `db:"sexo" json:"sexo"`
}

//ServidorController is used to export the API handler functions
type ServidorController struct{}

//GetServidor funtion returns the full list of "servidores" in the database
func (ctrl ServidorController) GetServidor(w http.ResponseWriter, r *http.Request) {
	q := `select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa)`
	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
		return
	}
	defer rows.Close()
	var servidores []Servidor
	var id, idpessoa, matriculainterna, siape int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ErrorBody{
				Reason: err.Error(),
			})
			log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
			return
		}
		date, _ := time.Parse("1969-02-12", datanascimento)
		servidores = append(servidores, Servidor{
			ID:                id,
			Siape:             siape,
			Idpessoa:          idpessoa,
			Nome:              nome,
			Matriculainterna:  matriculainterna,
			Nomeidentificacao: nomeidentificacao,
			Datanascimento:    date.Format("1969-02-12"),
			Sexo:              sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&servidores)
	log.Print("[MUX] " + " | 200 | " + r.Method + "  " + r.URL.Path)
	return
}

//GetServidorMat funtion returns the "servidor" matching a given id
func (ctrl ServidorController) GetServidorMat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mat := vars["matricula"] // URL parameter
	// Data security checking to be insterted here
	q := fmt.Sprintf(`select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
		inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa) where s.matricula_interna = %s`, mat) //String formating
	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	defer rows.Close()
	var servidores []Servidor
	var id, idpessoa, matriculainterna, siape int
	var nome, nomeidentificacao, datanascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &idpessoa, &matriculainterna, &nomeidentificacao, &nome, &datanascimento, &sexo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorBody{
				Reason: err.Error(),
			})
			log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
			return
		}
		date, _ := time.Parse("1969-02-12", datanascimento)
		servidores = append(servidores, Servidor{
			ID:                id,
			Siape:             siape,
			Idpessoa:          idpessoa,
			Nome:              nome,
			Matriculainterna:  matriculainterna,
			Nomeidentificacao: nomeidentificacao,
			Datanascimento:    date.Format("1969-02-12"),
			Sexo:              sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&servidores)
	log.Print("[MUX] " + " | 200 | " + r.Method + "  " + r.URL.Path)
	return
}

//PostServidor function reads a JSON body and store it in the database
func (ctrl ServidorController) PostServidor(w http.ResponseWriter, r *http.Request) {
	regexcheck := false
	var ser Servidor
	var Reasons []ErrorBody
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&ser)
	if errDecode != nil {
		log.Println(errDecode)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: errDecode.Error(),
		})
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	// REGEX CHEKING PHASE
	regex, _ := regexp.Compile(`^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)$`)
	if !regex.MatchString(ser.Datanascimento) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[data_nascimento] failed to match API requirements. It should look like this: 1969-02-12",
		})
	}
	regex, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !regex.MatchString(ser.Nome) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	regex, _ = regexp.Compile(`^([A-Z][a-z]+([ ]?[a-z]?['-]?[A-Z][a-z]+)*)$`)
	if !regex.MatchString(ser.Nomeidentificacao) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[nome_identificacao] failed to match API requirements. It should look like this: Firstname Middlename(optional) Lastname",
		})
	}
	regex, _ = regexp.Compile(`\b[MF]{1}\b`)
	if !regex.MatchString(ser.Sexo) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[sexo] failed to match API requirements. It should look like this: M for male, F for female",
		})
	}
	if regexcheck {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&Reasons)
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	// END OF REGEX CHEKING PHASE
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	b := md5.Sum([]byte(fmt.Sprintf(string(ser.Nome), string(timestamp))))
	bid := binary.BigEndian.Uint64(b[:])
	ser.Matriculainterna = int(bid % 9999999)
	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, siape, id_pessoa, matricula_interna, id_foto,
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, null, '%s', '%s');
			`, ser.Nome, ser.Nomeidentificacao, ser.Siape, ser.Idpessoa, ser.Matriculainterna,
		ser.Datanascimento, ser.Sexo) //String formating
	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		log.Print("[MUX] " + " | 500 | " + r.Method + "  " + r.URL.Path)
		return
	}
	defer rows.Close()
	w.Header().Add("location", r.URL.Host+"/api/servidor/"+strconv.Itoa(ser.Matriculainterna))
	w.WriteHeader(201)
	log.Print("[MUX] " + " | 201 | " + r.Method + "  " + r.URL.Path)
	return
}

func (ctrl ServidorController) Calculate(w http.ResponseWriter, r *http.Request) {
	var matrix [][]float64
	//matrixTwo := make([][]float64, 10)
	decoder := json.NewDecoder(r.Body)
	errDecode := decoder.Decode(&matrix)
	if errDecode != nil {
		log.Println(errDecode)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: errDecode.Error(),
		})
		log.Print("[MUX] " + " | 400 | " + r.Method + "  " + r.URL.Path)
		return
	}
	matrix = calc(matrix)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Result{
		Result: sum(matrix),
	})
}
func calc(matrix [][]float64) [][]float64 {
	for rowIndex, row := range matrix {
		relSum := 0.0
		for _, element := range row {
			relSum += math.Pow(element, 2)
		}
		relSum = relSum / float64(len(row))
		for index, element := range row {
			matrix[rowIndex][index] = math.Sqrt(element * relSum)
		}
	}
	return matrix
}
func sum(matrix [][]float64) float64 {
	relSum := 0.0
	for _, row := range matrix {
		for _, element := range row {
			relSum += element
		}
	}
	return relSum
}
