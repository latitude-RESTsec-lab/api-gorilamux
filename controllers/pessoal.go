package controllers

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pessoalAPI-gorilamux/db"
	"regexp"
	"time"

	"github.com/gorilla/mux"
)

type ErrorBody struct {
	Reason string `json:"reason"`
}

type Servidor struct {
	ID                 int    `db:"id, primarykey, autoincrement" json:"id"`
	Siape              int    `db:"siape" json:"siape"`
	Id_pessoa          int    `db:"id_pessoa" json:"id_pessoa"`
	Nome               string `db:"nome" json:"nome"`
	Matricula_interna  int    `db:"matricula_interna" json:"matricula_interna"`
	Nome_identificacao string `db:"nome_identificacao" json:"nome_identificacao"`
	Data_nascimento    string `db:"data_nascimento" json:"data_nascimento"`
	Sexo               string `db:"sexo" json:"sexo"`
}

type ServidorController struct{}

func (ctrl ServidorController) GetServidor(w http.ResponseWriter, r *http.Request) {
	q := `select s.id_servidor, s.siape, s.id_pessoa, s.matricula_interna, s.nome_identificacao,
		p.nome, p.data_nascimento, p.sexo from rh.servidor s
	inner join comum.pessoa p on (s.id_pessoa = p.id_pessoa)`

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		return
	}
	defer rows.Close()

	var pessoas []Servidor

	var id, id_pessoa, matricula_interna, siape int
	var nome, nome_identificacao, data_nascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &id_pessoa, &matricula_interna, &nome_identificacao, &nome, &data_nascimento, &sexo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorBody{
				Reason: err.Error(),
			})
			return
		}

		pessoas = append(pessoas, Servidor{
			ID:                 id,
			Siape:              siape,
			Id_pessoa:          id_pessoa,
			Nome:               nome,
			Matricula_interna:  matricula_interna,
			Nome_identificacao: nome_identificacao,
			Data_nascimento:    data_nascimento,
			Sexo:               sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&pessoas)

	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		return
	}

	return
}

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
		return
	}
	defer rows.Close()

	var pessoas []Servidor

	var id, id_pessoa, matricula_interna, siape int
	var nome, nome_identificacao, data_nascimento, sexo string
	for rows.Next() {
		err := rows.Scan(&id, &siape, &id_pessoa, &matricula_interna, &nome_identificacao, &nome, &data_nascimento, &sexo)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorBody{
				Reason: err.Error(),
			})
			return
		}

		pessoas = append(pessoas, Servidor{
			ID:                 id,
			Siape:              siape,
			Id_pessoa:          id_pessoa,
			Nome:               nome,
			Matricula_interna:  matricula_interna,
			Nome_identificacao: nome_identificacao,
			Data_nascimento:    data_nascimento,
			Sexo:               sexo,
		})
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(&pessoas)

	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		return
	}

	return
}

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
		return
	}

	// REGEX CHEKING PHASE
	regex, _ := regexp.Compile(`^(19[0-9]{2}|2[0-9]{3})-(0[1-9]|1[012])-([123]0|[012][1-9]|31)T([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])Z$`)
	if !regex.MatchString(ser.Data_nascimento) {
		regexcheck = true
		Reasons = append(Reasons, ErrorBody{
			Reason: "[data_nascimento] failed to match API requirements. It should look like this: 1969-02-12T00:00:00Z",
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
	if !regex.MatchString(ser.Nome_identificacao) {
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
	// regex, _ = regexp.Compile(`\b[0-9]+\b`)
	// if !regex.MatchString(strconv.Itoa(ser.Siape)) {
	// 	regexcheck = true
	// 	Reasons = append(Reasons, ErrorBody{
	// 		error_reason: "[siape] failed to match API requirements. It should be only numeric.",
	// 	})
	// }
	// regex, _ = regexp.Compile(`\b[0-9]+\b`)
	// if !regex.MatchString(strconv.Itoa(ser.Id_pessoa)) {
	// 	regexcheck = true
	// 	Reasons = append(Reasons, ErrorBody{
	// 		error_reason: "[id_pessoa] failed to match API requirements. It should be only numeric.",
	// 	})
	// }
	if regexcheck {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&Reasons)
		return
	}
	// END OF REGEX CHEKING PHASE
	timestamp := time.Now().Unix()
	b := md5.Sum([]byte(fmt.Sprintf(string(ser.Nome), string(timestamp))))
	bid := binary.BigEndian.Uint64(b[:])
	// log.Println(strconv.Atoi(string(b[:])))
	// ser.Data_nascimento = serTex.Data_nascimento
	// ser.ID, _ = strconv.Atoi(serTex.ID)

	q := fmt.Sprintf(`
		INSERT INTO rh.servidor_tmp(
			nome, nome_identificacao, siape, id_pessoa, matricula_interna, id_foto,
			data_nascimento, sexo)
			VALUES ('%s', '%s', %d, %d, %d, null, '%s', '%s');
			`, ser.Nome, ser.Nome_identificacao, ser.Siape, ser.Id_pessoa, bid%99999,
		ser.Data_nascimento, ser.Sexo) //String formating

	rows, err := db.GetDB().Query(q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorBody{
			Reason: err.Error(),
		})
		return
	}

	defer rows.Close()

	// var pessoas []Servidor
	w.WriteHeader(200)

	// if err != nil {
	// 	return
	// }

	return
}
