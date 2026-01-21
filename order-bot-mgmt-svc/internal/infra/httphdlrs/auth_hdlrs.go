package httphdlrs

import (
	"encoding/json"
	"net/http"

	"order-bot-mgmt-svc/internal/services/authsvc"
)

type Server interface {
	AuthService() *authsvc.Svc
}

type authRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func SignupHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		tokens, err := s.AuthService().Signup(req.Email, req.Password)
		if err != nil {
			switch err {
			case authsvc.ErrUserExists:
				http.Error(w, "user already exists", http.StatusConflict)
			case authsvc.ErrInvalidCredentials:
				http.Error(w, "invalid credentials", http.StatusBadRequest)
			default:
				http.Error(w, "failed to create user", http.StatusInternalServerError)
			}
			return
		}
		writeJSON(w, http.StatusCreated, tokens)
	}
}

func LoginHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		tokens, err := s.AuthService().Login(req.Email, req.Password)
		if err != nil {
			switch err {
			case authsvc.ErrInvalidCredentials:
				http.Error(w, "invalid credentials", http.StatusUnauthorized)
			default:
				http.Error(w, "failed to login", http.StatusInternalServerError)
			}
			return
		}
		writeJSON(w, http.StatusOK, tokens)
	}
}

func LogoutHandler(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := s.AuthService().Logout(req.RefreshToken); err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
