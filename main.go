package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
)

type credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

var creds credentials
var oconf *oauth2.Config

func init() {
	// Read credentials from creds.json
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	json.Unmarshal(file, &creds)

	//configure oAuth2 struct
	oconf = &oauth2.Config{
		ClientID:     creds.Cid,
		ClientSecret: creds.Csecret,
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello authers!!!")
}

func authGoogle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oconf.AuthCodeURL("safe"), http.StatusMovedPermanently)
}

func authGoogleCallback(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	code := queryValues.Get("code")
	log.Println("Google code: ", code)

	// convert code into a token
	tok, err := oconf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("Exchange error")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Client returns an authorized HTTP Client using the provided token
	client := oconf.Client(oauth2.NoContext, tok)

	// get the information using the http client
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	w.Write(data)
	return
}

func main() {
	log.Println("Startng goauth server: http://localhost:3000 ....")

	router := httprouter.New()

	chainHandlers := alice.New(context.ClearHandler, loggingHandler, recoveryHandler)

	router.GET("/", wrapHandler(chainHandlers.ThenFunc(index)))
	router.GET("/auth/google", wrapHandler(chainHandlers.ThenFunc(authGoogle)))
	router.GET("/auth/google/callback", wrapHandler(chainHandlers.ThenFunc(authGoogleCallback)))

	log.Fatal(http.ListenAndServe(":3000", router))
}
