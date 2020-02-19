package main

import "time"

type Eventos []Evento
type Events struct {
	Eventos Eventos `json:"events" bson:"events"`
}
type Evento struct {
	Event       string                   `json:"event" bson:"event"`
	Timestamp   time.Time                `json:"timestamp" bson:"timestamp"`
	Revenue     int                      `json:"revenue" bson:"revenue"`
	Custom_data []map[string]interface{} `json:"custom_data" bson:"custom_data"`
}
