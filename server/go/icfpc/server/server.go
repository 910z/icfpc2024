package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/uptrace/bun"
)

func New(db *bun.DB, port int) *Server {
	return &Server{
		port: port,
		db:   db,
	}
}

type Server struct {
	port int
	db   *bun.DB
}

func (s *Server) SetupEndpoints() {
	http.Handle("/", &app.Handler{})

	http.HandleFunc("/list", s.handleList)
}

func (s *Server) Run() error {
	slog.Info("Server listening", slog.Int("port", s.port))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil); err != nil {
		return err
	}

	return nil
}

func writeJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		writeError(w, http.StatusInternalServerError)

		return
	}
}

func writeError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
