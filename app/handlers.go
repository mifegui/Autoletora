package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
)

// Index retorna ao webserver o fato de que a api está rodando.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Rodando\n")
}

// EventoIndex retorna em json para o webserver todos os eventos coletados
// e armazenados na database até agora.
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

// EventoColetar recebe em json eventos que serão armazenados pela database
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

// TimelineFazer consome a api endpoint e agrupa os eventos do endpoint em uma timeline que é retornada para o webserver.
func TimelineFazer(w http.ResponseWriter, r *http.Request) {
	var endpoint Events                                                       // Struct que representa os dados do endpoint
	urlEndpoint := "http://storage.googleapis.com/dito-questions/events.json" // Url do endpoint
	response, err := http.Get(urlEndpoint)
	if err != nil {
		fmt.Printf("A solicitação http falhou com o erro: %s\n", err)
		panic(err)
	} else {
		data, err := ioutil.ReadAll(response.Body) // armazena os dados que queremos em data
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(data, &endpoint); err != nil { // Se não conseguir codificar o json para a struct...
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422)                                     // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil { // Mostra o erro
				panic(err)
			}
		}
	}

	// Agora a struct endpoint está preenchida com os dados do .json
	// Passaremos pra struct formatada
	m := make(map[string]Transaction)         // mapa cuja key é a string de transaction_id e o valor é a própria struct transaction
	for _, evento := range endpoint.Eventos { // Para cada evento
		var trans Transaction // construiremos uma transação
		trans.Timestamp = evento.Timestamp
		trans.Revenue = evento.Revenue
		var prod Product
		for _, csdata := range evento.Custom_data { // Para cada key e valor de cada evento
			if csdata["key"] == "transaction_id" {
				trans.Transaction_id = csdata["value"].(string)
			}
			if csdata["key"] == "store_name" {
				trans.Store_name = csdata["value"].(string)
			}
			if csdata["key"] == "product_name" {
				prod.Name = csdata["value"].(string)
			}
			if csdata["key"] == "product_price" {
				prod.Price = int(csdata["value"].(float64))
			}

		} // processamos o evento e as keys e values, agora vamos juntar com a transação de mesmo id se existir
		tpassada, existe := m[trans.Transaction_id] // Checa se o a transid do evento atual já foi vista anteriormente
		if existe {                                 // booleana
			// Passa os valores da trans atual para a trans passada
			if evento.Event == "comprou" { // só existe um comprou por transação
				tpassada.Revenue = trans.Revenue
				tpassada.Store_name = trans.Store_name
			} else if evento.Event == "comprou-produto" {
				tpassada.Products = append(tpassada.Products, prod)
			}
			m[trans.Transaction_id] = tpassada // troca a antiga pela modificada
		} else {

			// Coloca a trans atual no mapa
			trans.Products = append(trans.Products, prod)
			m[trans.Transaction_id] = trans

		}

	} // passamos por todos os eventos e agora temos as transações em m, colocaremos elas ordenadamente na timeline
	var timelined Timeline
	for _, trans := range m { // para cada transação
		timelined.Trans = append(timelined.Trans, trans) // coloque-a no vetor de transações de timelined
	}
	sort.Sort(Transactions(timelined.Trans)) // Coloque em ordem alfebética as Transactions de timelined, como especificado em timeline.go

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422)                                           // unprocessable entity
	if err := json.NewEncoder(w).Encode(timelined); err != nil { // Manda o timelined
		panic(err)
	}

}
