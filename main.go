package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
)

func index(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintf(w, "Hello authers!!!")
	return &appError{nil, "error", 500}
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Test Handler!!!")
}

func main() {

	router := httprouter.New()

	chainHandlers := alice.New(context.ClearHandler, loggingHandler, recoveryHandler)

	router.GET("/", wrapHandler(chainHandlers.ThenFunc(index)))
	router.GET("/test", wrapHandler(chainHandlers.ThenFunc(test)))

	log.Fatal(http.ListenAndServe(":3000", router))
}
