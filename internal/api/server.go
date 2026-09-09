package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/JMar2021/sports-data-platform/internal/repository"
)

type Server struct {
	Repository *repository.Repository
	Mux        *http.ServeMux
	HTTPServer *http.Server
	// Add any necessary fields for the server, such as configuration, logger, etc.
}

func NewServer(repo *repository.Repository, httpAddr string) *Server {
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	server := &Server{
		Repository: repo,
		Mux:        mux,
		HTTPServer: httpServer,
	}
	mux.HandleFunc("/games", server.handleGames)
	mux.HandleFunc("/health", healthHandler)
	return server
}

func (s *Server) handleGames(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	location := now.Location()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		location,
	)

	start := today.AddDate(0, 0, -1)
	end := today.AddDate(0, 0, 8)

	ctx := r.Context()
	games, err := s.Repository.GetGames(
		ctx,
		start,
		end,
	)
	if err != nil {
		http.Error(w, "failed to get games", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func (s *Server) Start() error {
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}
