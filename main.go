package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Go Auth Router... ")
}

func main() {
	fmt.Println("Startng goauth server: localhost:8080 ....")

	router := httprouter.New()
	router.GET("/", index)

	log.Fatal(http.ListenAndServe(":8080", router))

}
