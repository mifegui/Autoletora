package main

import (
	"log"
	"net/http"
	"time"
)

// Logger é um wrapper para os outros handlers, calcula o tempo que demorou e registra mais alguns parâmetros que o handler usa.
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
