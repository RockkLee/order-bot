package httphdlrs

import (
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/util/errutil"
	"order-bot-mgmt-svc/internal/util/validatorutil"

	"order-bot-mgmt-svc/internal/services/authsvc"
)

type AuthServer interface {
	AuthService() *authsvc.Svc
}

const AuthPrefix = "/auth"

func AuthHdlr(s AuthServer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", signupHdlrFunc(s))
	mux.HandleFunc("POST /login", loginHdlrFunc(s))
	mux.HandleFunc("POST /logout", logoutHldrFunc(s))
	return mux
}

func signupHdlrFunc(s AuthServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJsonRequest[signupRequest](w, r)
		if !ok {
			return
		}
		if err := validatorutil.RequiredStrings(req); err != nil {
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
			return
		}
		tokens, err := s.AuthService().Signup(req.Email, req.Password)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			switch err {
			case authsvc.ErrUserExists:
				http.Error(w, authsvc.ErrUserExists.Error(), http.StatusConflict)
			case authsvc.ErrInvalidCredentials:
				http.Error(w, authsvc.ErrInvalidCredentials.Error(), http.StatusBadRequest)
			default:
				http.Error(w, "failed to create user", http.StatusInternalServerError)
			}
			return
		}
		writeJSON(w, http.StatusCreated, tokens)
	}
}

func loginHdlrFunc(s AuthServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJsonRequest[loginRequest](w, r)
		if !ok {
			return
		}
		if err := validatorutil.RequiredStrings(req); err != nil {
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
			return
		}
		tokens, err := s.AuthService().Login(req.Email, req.Password)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			switch err {
			case authsvc.ErrInvalidCredentials:
				http.Error(w, authsvc.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
			default:
				http.Error(w, "failed to login", http.StatusInternalServerError)
			}
			return
		}
		writeJSON(w, http.StatusOK, tokens)
	}
}

func logoutHldrFunc(s AuthServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := decodeJsonRequest[logoutRequest](w, r)
		if !ok {
			return
		}
		if err := validatorutil.RequiredStrings(req); err != nil {
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
			return
		}
		if err := s.AuthService().Logout(req.RefreshToken); err != nil {
			slog.Error(errutil.FormatErrChain(err))
			WriteError(w, http.StatusUnauthorized, authsvc.ErrInvalidRefreshToken.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": authsvc.ErrLoggedOut.Error()})
	}
}
