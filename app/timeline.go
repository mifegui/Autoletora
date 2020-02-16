// Structs para serem usadas em TimelineFazer em handlers.go
package main

import "time"

type Events struct {
	Eventos []Evento `json:"events"`
}
type Timeline struct {
	Trans []Transaction
}

type Transaction struct {
	Timestamp      time.Time `json:"timestamp"`
	Revenue        int       `json:"revenue"`
	Transaction_id int       `json:"transaction_id"`
	Store_name     string    `json:"store_name"`
	Products       []Product `json:"products"`
}

type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}
