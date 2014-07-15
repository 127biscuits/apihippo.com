package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/127biscuits/apihippo.com/api"
	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Panic("I can't start the service without finding the template folder!")
	}
	t.Execute(w, nil)
}

func main() {
	// TODO: add a config file
	port := 8000

	idRegExp := "/{id:[a-f0-9]{24}}"

	r := mux.NewRouter()

	// This regexp looks kinda hacky, but I don't mind about the rest of the
	// host.
	// It needs to be here because we want the host matching first.
	s := r.Host("cdn.{_:.*}").Subrouter()
	s.HandleFunc(idRegExp, api.FakeCDNHandler)

	r.HandleFunc("/", api.GetHandler).
		Methods("GET").Headers("Accept", "application/json")
	r.HandleFunc(idRegExp, api.GetHippoHandler).
		Methods("GET").Headers("Accept", "application/json")
	r.HandleFunc("/", api.PostHandler).
		Methods("POST")
	r.HandleFunc(idRegExp, api.PutHippoHandler).
		Methods("PUT")

	r.HandleFunc("/", indexHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Print("Listening at port ", port)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
