package httphdlrs

import (
	"encoding/json"
	"net/http"

	"order-bot-mgmt-svc/internal/services"
)

type authRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	tokens, err := s.authService().Signup(req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrUserExists:
			http.Error(w, "user already exists", http.StatusConflict)
		case services.ErrInvalidCredentials:
			http.Error(w, "invalid credentials", http.StatusBadRequest)
		default:
			http.Error(w, "failed to create user", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusCreated, tokens)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	tokens, err := s.authService().Login(req.Email, req.Password)
	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "failed to login", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, tokens)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := s.authService().Logout(req.RefreshToken); err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
