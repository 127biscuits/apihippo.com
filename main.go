package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/127biscuits/apihippo.com/api"
	"github.com/127biscuits/apihippo.com/settings"
	"github.com/gorilla/mux"
)

var settingsPath string

func init() {
	flag.StringVar(&settingsPath, "s", "settings.json", "JSON configuration")
}

func main() {
	flag.Parse()

	// TODO: get it from a path
	if err := settings.Load(settingsPath); err != nil {
		log.Fatal("Error reading config file: ", err)
	}

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

	log.Print("Listening at port ", settings.Config.Port)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", settings.Config.Port), r))
}
