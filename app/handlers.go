package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

func TimelineFazer(w http.ResponseWriter, r *http.Request) {
	var endpoint Events                                                       // Struct que representa os dados do endpoint
	urlEndpoint := "http://storage.googleapis.com/dito-questions/events.json" // Url do endpoint
	response, err := http.Get(urlEndpoint)

	if err != nil {
		fmt.Printf("A solicitação http falhou com o erro: %s\n", err)
		panic(err)
	} else {
		data, err := ioutil.ReadAll(response.Body)
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
	fmt.Printf("The struct returned before marshalling\n\n")
	fmt.Printf("%+v\n\n\n\n", endpoint)

	// The MarshalIndent function only serves to pretty print, json.Marshal() is what would normally be used
	byteArray, err := json.MarshalIndent(endpoint, "", "  ")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("The JSON response returned when the struct is marshalled\n\n")
	fmt.Println(string(byteArray))
	// Agora a struct endpoint está preenchida com os dados do .json
	// Passaremos pra struct formatada
	m := make(map[int]Transaction)            // mapa cuja key é o int de transaction_id e o valor é a própria struct transaction
	for _, evento := range endpoint.Eventos { // Para cada evento
		var trans Transaction
		trans.Timestamp = evento.Timestamp
		if evento.Event == "comprou" {
			trans.Revenue = evento.Revenue
		}
		var prod Product
		for _, csdata := range evento.Custom_data { // Para cada key e valor de cada evento
			if csdata["Key"] == "transaction_id" {
				trans.Transaction_id = csdata["Value"].(int)
			}
			if evento.Event == "comprou" {
				if csdata["Key"] == "store_name" {
					trans.Store_name = csdata["Value"].(string)
				}
			} else if evento.Event == "comprou-produto" {
				if csdata["Key"] == "product_name" {
					prod.Name = csdata["Value"].(string)
				} else if csdata["Key"] == "product_price" {
					prod.Price = csdata["Value"].(int)
				}
			}

		} // processamos o evento e as keys e values, agora vamos juntar com a transação de mesmo id se existir
		tpassada, existe := m[trans.Transaction_id] // Checa se o a transid do evento atual já foi vista anteiormente
		if existe {                                 // booleana
			// Passa os valores da trans atual para a trans passada
			if evento.Event == "comprou" { // só existe um comprou por transação
				tpassada.Revenue = tpassada.Revenue
				tpassada.Store_name = tpassada.Store_name
			} else if evento.Event == "comprou-produto" {
				tpassada.Products = append(tpassada.Products, prod)
			}
			m[trans.Transaction_id] = tpassada // troca a antiga pela modificada
		} else {
			// Coloca a trans atual no mapa
			m[trans.Transaction_id] = trans
		}

	} // passamos por todos os eventos e agora temos as transações em m, colocaremos elas ordenadamente na timeline
	var timelined Timeline
	var maisNova Transaction
	for len(m) != 0 { // enquanto o mapa não estiver vazio
		for _, trans := range m {
			if maisNova.Transaction_id == 0 { // Se o maisnova estiver vazio
				maisNova = trans
			} else {
				if trans.Timestamp.After(maisNova.Timestamp) { // o atual é mais novo que o que estamos guardando
					maisNova = trans
				}
			}

		}
		timelined.Trans = append(timelined.Trans, maisNova)
		delete(m, maisNova.Transaction_id) // deleta o maisNova do mapa
	} // agora o timelined tem os Transactions em ordem temporal

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(422)                                           // unprocessable entity
	if err := json.NewEncoder(w).Encode(timelined); err != nil { // Manda o timelined
		panic(err)
	}

}
