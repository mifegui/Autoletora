// Structs para serem usadas em TimelineFazer em handlers.go
package main

import "time"

type Transactions []Transaction

type Timeline struct {
	Trans Transactions `json:"timeline"`
}

type Transaction struct {
	Timestamp      time.Time `json:"timestamp"`
	Revenue        int       `json:"revenue"`
	Transaction_id string    `json:"transaction_id"`
	Store_name     string    `json:"store_name"`
	Products       []Product `json:"products"`
}
type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Implementações da Interface Sort
func (t Transactions) Len() int {
	return len(t)
}

func (t Transactions) Less(i, j int) bool {
	return t[i].Timestamp.After(t[j].Timestamp)
}

func (t Transactions) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
