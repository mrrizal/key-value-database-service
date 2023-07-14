package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrrizal/key-value-database/handler"
	"github.com/mrrizal/key-value-database/service"
)

func main() {
	storeService := service.NewStoreService()
	storeHandler := handler.NewStoreHandler(storeService)

	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", storeHandler.Put).Methods("PUT")
	r.HandleFunc("/v1/{key}", storeHandler.Get).Methods("GET")
	r.HandleFunc("/v1/{key}", storeHandler.Delete).Methods("DELETE")

	log.Println("starting the server...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
