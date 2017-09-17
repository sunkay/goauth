package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("Error: ", err)
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

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	type profile struct {
		First string
		Last  string
	}

	p := profile{"sun", "kay"}
	js, err := json.Marshal(p)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("json:", js)

	w.Write(js)
}

func authGoogle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, oconf.AuthCodeURL("safe"), http.StatusMovedPermanently)
}

func authGoogleCallback(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// retrieve query param code
	queryValues := r.URL.Query()
	code := queryValues.Get("code")
	fmt.Println("Google code: ", code)

	// convert code into a token
	tok, _ := oconf.Exchange(oauth2.NoContext, code)

	// Client returns an authorized HTTP Client using the provided token
	client := oconf.Client(oauth2.NoContext, tok)

	// get the information using the http client
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	w.Write(data)
}

func main() {
	fmt.Println("Startng goauth server: localhost:8080 ....")

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/auth/google", authGoogle)
	router.GET("/auth/google/callback", authGoogleCallback)

	log.Fatal(http.ListenAndServe(":3001", handlers.LoggingHandler(os.Stdout, router)))

}
