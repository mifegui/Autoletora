package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"strings"
	"time"
)

var client *mongo.Client
var mongodbName string = "mongodb://localhost:27017"

func init() {
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

}

func getEventos(w http.ResponseWriter) (Events, error) {
	var evts Events
	collection := client.Database("autoletora").Collection("eventos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return evts, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var e Evento
		cursor.Decode(&e)
		evts.Eventos = append(evts.Eventos, e)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return evts, err
	}
	return evts, nil
}

func getMatch(w http.ResponseWriter, input string) (Completed, error) {

	// Procurar matchs na database
	var mts Completed
	collection := client.Database("autoletora").Collection("eventos")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{}) // Procura pela database inteira
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return mts, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var e Evento
		cursor.Decode(&e)
		if strings.HasPrefix(e.Event, input) {
			existe := false
			for _, m := range mts.Matchs {
				if m == e.Event {
					existe = true
					break
				}
			}
			if existe == false {
				mts.Matchs = append(mts.Matchs, e.Event)
			}
		}
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return mts, err
	}
	return mts, nil

}

func createEvento(e Evento) *mongo.InsertOneResult {
	collection := client.Database("autoletora").Collection("eventos")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, e)
	return result

}
