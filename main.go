package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func authGoogle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, oconf.AuthCodeURL("safe"), http.StatusMovedPermanently)
}

func authGoogleCallback(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// retrieve query param code
	queryValues := r.URL.Query()
	code := queryValues.Get("code")
	log.Println("Google code: ", code)

	// convert code into a token
	tok, err := oconf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	// Client returns an authorized HTTP Client using the provided token
	client := oconf.Client(oauth2.NoContext, tok)

	// get the information using the http client
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	w.Write(data)
}

func main() {
	log.Println("Startng goauth server: localhost:8080 ....")

	router := httprouter.New()
	//router.GET("/", index)
	router.GET("/auth/google", authGoogle)
	router.GET("/auth/google/callback", authGoogleCallback)

	log.Fatal(http.ListenAndServe(":3001", handlers.LoggingHandler(os.Stdout, router)))

}
