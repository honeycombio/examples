package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HoneycombMiddleware(mainHandler))
	r.HandleFunc("/signup", HoneycombMiddleware(signupHandler))
	r.HandleFunc("/login", HoneycombMiddleware(loginHandler))
	r.HandleFunc("/logout", HoneycombMiddleware(logoutHandler))
	r.HandleFunc("/shout", HoneycombMiddleware(shoutHandler))
	log.Print("Serving app on localhost:8888 ....")
	log.Fatal(http.ListenAndServe(":8888", r))
}
