package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/httphdlrs"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	// Register routes
	mux.HandleFunc("/", s.helloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.Handle(
		fmt.Sprintf("%s/", httphdlrs.AuthPrefix),
		http.StripPrefix(httphdlrs.AuthPrefix, httphdlrs.AuthHdlr(s)),
	)

	// Wrap the mux with CORS middleware
	middlewareStack := createMiddlewareStack(
		corsMiddleware,
	)
	return middlewareStack(mux)
}

func (s *Server) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := s.dbService().Health()
	if err != nil {
		http.Error(w, "Failed to check database health", http.StatusServiceUnavailable)
		return
	}
	resp, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
