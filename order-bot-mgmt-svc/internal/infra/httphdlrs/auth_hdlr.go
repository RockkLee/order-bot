package httphdlrs

import (
	"encoding/json"
	"errors"
	"net/http"

	"order-bot-mgmt-svc/internal/services/authsvc"
)

type Server interface {
	AuthService() *authsvc.Svc
}

const AuthPrefix = "/auth"

func AuthHdlr(s Server) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", signupHdlrFunc(s))
	mux.HandleFunc("POST /login", loginHdlrFunc(s))
	mux.HandleFunc("POST /logout", logoutHldrFunc(s))
	return mux
}

func signupHdlrFunc(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeAuthRequest(w, r)
		if !ok {
			return
		}
		tokens, err := s.AuthService().Signup(req.Email, req.Password)
		if err != nil {
			handleSignupError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, tokens)
	}
}

func loginHdlrFunc(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeAuthRequest(w, r)
		if !ok {
			return
		}
		tokens, err := s.AuthService().Login(req.Email, req.Password)
		if err != nil {
			handleLoginError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, tokens)
	}
}

func logoutHldrFunc(s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeAuthRequest(w, r)
		if !ok {
			return
		}
		if err := s.AuthService().Logout(req.RefreshToken); err != nil {
			WriteError(w, http.StatusUnauthorized, ErrMsgInvalidRefreshToken)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": ErrMsgLoggedOut})
	}
}

func decodeAuthRequest(w http.ResponseWriter, r *http.Request) (authRequest, bool) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
		return authRequest{}, false
	}
	return req, true
}

func handleSignupError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authsvc.ErrUserExists):
		WriteError(w, http.StatusConflict, ErrMsgUserAlreadyExists)
	case errors.Is(err, authsvc.ErrInvalidCredentials):
		WriteError(w, http.StatusBadRequest, ErrMsgInvalidCredentials)
	default:
		WriteError(w, http.StatusInternalServerError, ErrMsgFailedCreateUser)
	}
}

func handleLoginError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, authsvc.ErrInvalidCredentials):
		WriteError(w, http.StatusUnauthorized, ErrMsgInvalidCredentials)
	default:
		WriteError(w, http.StatusInternalServerError, ErrMsgFailedLogin)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
