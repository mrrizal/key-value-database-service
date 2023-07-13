package handler

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrrizal/key-value-database/service"
)

type StoreHandler struct {
	svc service.StoreService
}

func NewStoreHandler(svc service.StoreService) *StoreHandler {
	return &StoreHandler{
		svc: svc,
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
