package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// hanling JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// Error
type apiError struct {
	Error string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

// makeHTTPHandleFunc is decorator to http.HandlerFunc
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	log.Println("API server runing on port: ", s.listenAddr)
	router := mux.NewRouter()
	router.HandleFunc("/health", makeHTTPHandleFunc(s.handleHealth))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		s.handleGetAccount(w, r)
	case "POST":
		s.handleCreateAccount(w, r)
	case "DELETE":
		s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("%s Method not allowed", r.Method)
	}
	return nil
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		WriteJSON(w, http.StatusOK, "service is running")
	} else {
		return fmt.Errorf("%s Method not allowed", r.Method)
	}
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	// account := NewAccount("Alok", "Tripathi")
	return WriteJSON(w, http.StatusOK, vars)

}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
