package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

var client *mongo.Client
var mongodbName string = "mongodb://localhost:27017"

func main() {

	log.Println("Iniciando/Conectando database...")
	auxclient, err := mongo.NewClient(options.Client().ApplyURI(mongodbName))
	if err != nil {
		log.Fatal(err)
	}
	client = auxclient
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongodbName))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Não conectado à database: ", err)
	}

	log.Println("Iniciando servidor web")
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
