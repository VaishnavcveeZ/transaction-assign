package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/transaction", CreateTransaction).Methods(http.MethodPost)
	router.HandleFunc("/transaction", DeleteTransaction).Methods(http.MethodDelete)

	router.HandleFunc("/statistics", GetStatics).Methods(http.MethodGet)

	router.HandleFunc("/location", SetUserCity).Methods(http.MethodPost)
	router.HandleFunc("/location", ResetUserCity).Methods(http.MethodPut)

	log.Fatal(http.ListenAndServe(":4000", router))
}
