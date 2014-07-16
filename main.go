package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/127biscuits/apihippo.com/api"
	"github.com/gorilla/mux"
)

func main() {
	// TODO: add a config file
	port := 8000

	idRegExp := "/{id:[a-f0-9]{24}}"

	r := mux.NewRouter()

	// This regexp looks kinda hacky, but I don't mind about the rest of the
	// host.
	// It needs to be here because we want the host matching first.
	cdn := r.Host("cdn.{_:.*}").Subrouter()
	cdn.HandleFunc(idRegExp, api.FakeCDNHandler)

	random := r.Host("random.{_:.*}").Subrouter()
	random.HandleFunc("/", api.RandomHippoHandler)

	r.HandleFunc("/", api.GetHandler).
		Methods("GET").Headers("Accept", "application/json")
	r.HandleFunc(idRegExp, api.GetHippoHandler).
		Methods("GET").Headers("Accept", "application/json")
	r.HandleFunc("/", api.PostHandler).
		Methods("POST")
	r.HandleFunc(idRegExp, api.PutHippoHandler).
		Methods("PUT")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Print("Listening at port ", port)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
