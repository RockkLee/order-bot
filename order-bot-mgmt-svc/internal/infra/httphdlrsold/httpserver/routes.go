package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-bot-mgmt-svc/internal/infra/httphdlrsold"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	// Register routes
	mux.HandleFunc("/", s.helloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.Handle(
		fmt.Sprintf("%s/", httphdlrsold.AuthPrefix),
		http.StripPrefix(httphdlrsold.AuthPrefix, httphdlrsold.AuthHdlr(s)),
	)
	mux.Handle(
		fmt.Sprintf("%s/", httphdlrsold.MenuPrefix),
		http.StripPrefix(httphdlrsold.MenuPrefix, httphdlrsold.MenuHdlr(s)),
	)
	mux.Handle(
		fmt.Sprintf("%s/", httphdlrsold.BotPrefix),
		http.StripPrefix(httphdlrsold.BotPrefix, httphdlrsold.BotHdlr(s)),
	)

	// Wrap the mux with CORS middleware
	middlewareStack := createMiddlewareStack(
		corsMiddleware(s),
		authMiddleware(s),
	)
	return middlewareStack(mux)
}

func (s *Server) helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		httphdlrsold.WriteError(w, http.StatusInternalServerError, httphdlrsold.ErrMsgFailedMarshalResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("%s: %v", httphdlrsold.LogMsgFailedWriteResponse, err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := s.dbService().Health()
	if err != nil {
		httphdlrsold.WriteError(w, http.StatusServiceUnavailable, httphdlrsold.ErrMsgFailedCheckDatabaseHealth)
		return
	}
	resp, err := json.Marshal(stats)
	if err != nil {
		httphdlrsold.WriteError(w, http.StatusInternalServerError, httphdlrsold.ErrMsgFailedMarshalHealthCheck)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("%s: %v", httphdlrsold.LogMsgFailedWriteResponse, err)
	}
}
