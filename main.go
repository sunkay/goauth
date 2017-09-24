package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello authers!!!")
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Test Handler!!!")
}

func main() {
	chainHandlers := alice.New(context.ClearHandler, loggingHandler, recoveryHandler)

	http.Handle("/", chainHandlers.ThenFunc(index))
	http.Handle("/test", chainHandlers.ThenFunc(test))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
