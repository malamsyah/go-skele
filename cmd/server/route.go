package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) registerRoutes() http.Handler {
	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/ping").HandlerFunc(Ping)
	return router
}

func Ping(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
	})
}
