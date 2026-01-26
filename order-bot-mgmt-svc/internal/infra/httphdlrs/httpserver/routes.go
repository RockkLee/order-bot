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
	mux.Handle(
		fmt.Sprintf("%s/", httphdlrs.MenuPrefix),
		http.StripPrefix(httphdlrs.MenuPrefix, httphdlrs.MenuHdlr(s)),
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
		httphdlrs.WriteError(w, http.StatusInternalServerError, httphdlrs.ErrMsgFailedMarshalResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("%s: %v", httphdlrs.LogMsgFailedWriteResponse, err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := s.dbService().Health()
	if err != nil {
		httphdlrs.WriteError(w, http.StatusServiceUnavailable, httphdlrs.ErrMsgFailedCheckDatabaseHealth)
		return
	}
	resp, err := json.Marshal(stats)
	if err != nil {
		httphdlrs.WriteError(w, http.StatusInternalServerError, httphdlrs.ErrMsgFailedMarshalHealthCheck)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("%s: %v", httphdlrs.LogMsgFailedWriteResponse, err)
	}
}
