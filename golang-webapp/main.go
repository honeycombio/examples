package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HoneycombMiddleware(mainHandler)).Methods("GET")

	r.HandleFunc("/signup", HoneycombMiddleware(signupHandlerGet)).Methods("GET")
	r.HandleFunc("/signup", HoneycombMiddleware(signupHandlerPost)).Methods("POST")

	r.HandleFunc("/login", HoneycombMiddleware(loginHandlerGet)).Methods("GET")
	r.HandleFunc("/login", HoneycombMiddleware(loginHandlerPost)).Methods("POST")

	r.HandleFunc("/logout", HoneycombMiddleware(logoutHandler)).Methods("POST")
	r.HandleFunc("/shout", HoneycombMiddleware(shoutHandler)).Methods("POST")

	log.Print("Serving app on localhost:8888 ....")
	log.Fatal(http.ListenAndServe(":8888", r))
}
