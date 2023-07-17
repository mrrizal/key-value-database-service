package handler

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrrizal/key-value-database/logger"
	"github.com/mrrizal/key-value-database/service"
)

type StoreHandler struct {
	svc    service.StoreService
	logger logger.TransactionLogger
}

func NewStoreHandler(svc service.StoreService, logger logger.TransactionLogger) *StoreHandler {
	return &StoreHandler{
		svc:    svc,
		logger: logger,
	}
}

func (s *StoreHandler) Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}()

	if err := s.svc.Put(key, string(value)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated)
}

func (s *StoreHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := s.svc.Get(key)
	if err != nil {
		if err == service.ErrorNoSuchKey {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func (s *StoreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if err := s.svc.Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}
