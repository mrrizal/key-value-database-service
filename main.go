package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/mrrizal/key-value-database/handler"
	"github.com/mrrizal/key-value-database/logger"
	"github.com/mrrizal/key-value-database/service"
)

func main() {
	logger, err := logger.InitializeTransactionLog("transaction.log")
	if err != nil {
		log.Fatal(err.Error())
	}

	storeService := service.NewStoreService()
	storeHandler := handler.NewStoreHandler(storeService, logger)

	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", storeHandler.Put).Methods("PUT")
	r.HandleFunc("/v1/{key}", storeHandler.Get).Methods("GET")
	r.HandleFunc("/v1/{key}", storeHandler.Delete).Methods("DELETE")

	go func() {
		log.Println("starting the server...")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-signalChannel:
			log.Println("closing event channel...")
			logger.Close()
			return
		default:
			logger.ReadEvents()
		}
	}
}
