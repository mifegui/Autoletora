package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//"strconv"
	//"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Rodando\n")
}

func EventoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(eventos); err != nil {
		panic(err)
	}
}

//func TodoShow(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	var todoId int
//	var err error
//	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
//		panic(err)
//	}
//	todo := RepoFindTodo(todoId)
//	if todo.Id > 0 {
//		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//		w.WriteHeader(http.StatusOK)
//		if err := json.NewEncoder(w).Encode(todo); err != nil {
//			panic(err)
//		}
//		return
//	}
//
//	// If we didn't find it, 404
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusNotFound)
//	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
//		panic(err)
//	}
//
//}

/*
Teste com esse comando:

curl -H "Content-Type: application/json" -d '{"event":"buy", "timestamp":"2016-09-22T13:57:31.2311892-04:00"}' http://localhost:8080/coletar
depois veja em localhost:8080/eventos

*/
func EventoColetar(w http.ResponseWriter, r *http.Request) {
	var evento Evento
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576)) // 1 Mega de limite para leitura
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &evento); err != nil { // Se não conseguir codificar o json para a struct...
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)                                     // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil { // Mostra o erro
			panic(err)
		}
	}

	e := RepoCriarEvento(evento) // Guarda na database
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		panic(err)
	}
}
