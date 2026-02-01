package httphdlrs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/errutil"
	"order-bot-mgmt-svc/internal/util/validatorutil"

	"order-bot-mgmt-svc/internal/services/authsvc"
)

type AuthServer interface {
	AuthService() *authsvc.Svc
	BotService() *botsvc.Svc
	WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error
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
		var (
			tokens models.TokenPair
			userId string
		)
		err := s.WithTx(r.Context(), func(ctx context.Context, tx store.Tx) error {
			var err error
			tokens, userId, err = s.AuthService().Signup(ctx, tx, req.Email, req.Password)
			if err != nil {
				return err
			}
			if errBot := s.BotService().CreateBot(ctx, tx, req.BotName, userId); errBot != nil {
				return fmt.Errorf("httphdlrs.CreateBot: %w", errBot)
			}
			return nil
		})
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			switch {
			case errors.Is(err, authsvc.ErrUserExists):
				http.Error(w, authsvc.ErrUserExists.Error(), http.StatusConflict)
			case errors.Is(err, authsvc.ErrInvalidCredentials):
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
		tokens, err := s.AuthService().Login(r.Context(), req.Email, req.Password)
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
		if err := s.AuthService().Logout(r.Context(), req.RefreshToken); err != nil {
			slog.Error(errutil.FormatErrChain(err))
			WriteError(w, http.StatusUnauthorized, authsvc.ErrInvalidRefreshToken.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": authsvc.ErrLoggedOut.Error()})
	}
}
