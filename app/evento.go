package main

import "time"

type Evento struct {
	Id          int
	Event       string                   `json:"event"`
	Timestamp   time.Time                `json:"timestamp"`
	Revenue     int                      `json:"revenue"`
	Custom_data []map[string]interface{} `json:"custom_data,omitempty"`
}

//type Data struct {
//	Key    string `json:"key"`
//	iValue int    `json:"value"`
//	Value  string `json:"value"`
//}

type Eventos []Evento
