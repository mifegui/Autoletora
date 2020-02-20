package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Iniciando servidor web")
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
