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

	r := mux.NewRouter()

	r.HandleFunc("/", api.GetHandler).Methods("GET").Headers("Accept", "application/json")
	r.HandleFunc("/{id}", api.GetHippoHandler).Methods("GET").Headers("Accept", "application/json")
	r.HandleFunc("/", api.PostHandler).Methods("POST")

	r.HandleFunc("/", indexHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Print("Listening at port ", port)
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
