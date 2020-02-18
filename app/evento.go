package main

import "time"

type Eventos []Evento
type Events struct {
	Eventos Eventos `json:"events"`
}
type Evento struct {
	Id          int
	Event       string                   `json:"event"`
	Timestamp   time.Time                `json:"timestamp"`
	Revenue     int                      `json:"revenue"`
	Custom_data []map[string]interface{} `json:"custom_data"`
}
