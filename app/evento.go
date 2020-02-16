package main

import "time"

type Evento struct {
	Event       string    `json:"event"`
	Timestamp   time.Time `json:"timestamp"`
	Revenue     int       `json:"revenue"`
	Custom_data []Data    `json:"custom_data"`
	Id          int
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Eventos []Evento
