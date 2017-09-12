package main

import (
	"log"
	"net/http"
)

func httpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
