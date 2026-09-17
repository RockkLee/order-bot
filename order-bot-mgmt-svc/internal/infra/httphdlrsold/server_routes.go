package httphdlrsold

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	// Register routes
	mux.HandleFunc("/", s.helloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.Handle(
		fmt.Sprintf("%s/", AuthPrefix),
		http.StripPrefix(AuthPrefix, AuthHdlr(s)),
	)
	mux.Handle(
		fmt.Sprintf("%s/", MenuPrefix),
		http.StripPrefix(MenuPrefix, MenuHdlr(s)),
	)
	mux.Handle(
		fmt.Sprintf("%s/", BotPrefix),
		http.StripPrefix(BotPrefix, BotHdlr(s)),
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
		WriteError(w, http.StatusInternalServerError, ErrMsgFailedMarshalResponse.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("%s: %v", LogMsgFailedWriteResponse, err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db().Health()
	if err != nil {
		WriteError(w, http.StatusServiceUnavailable, ErrMsgFailedCheckDatabaseHealth.Error())
		return
	}
	resp, err := json.Marshal(stats)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, ErrMsgFailedMarshalHealthCheck.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("%s: %v", LogMsgFailedWriteResponse, err)
	}
}
