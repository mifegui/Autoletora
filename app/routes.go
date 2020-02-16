package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"EventoIndex",
		"GET",
		"/eventos",
		EventoIndex,
	},
	Route{
		"EventoColetar",
		"POST",
		"/coletar",
		EventoColetar,
	},
	//	Route{
	//		"TodoShow",
	//		"GET",
	//		"/todos/{todoId}",
	//		TodoShow,
	//	},
}
