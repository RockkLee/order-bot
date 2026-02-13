package httphdlr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/errutil"

	"github.com/gin-gonic/gin"
)

type AuthServer interface {
	AuthService() *authsvc.Svc
	BotService() *botsvc.Svc
	WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error
}

const AuthPrefix = "/auth"

func RegisterAuthRoutes(r gin.IRoutes, s AuthServer) {
	r.POST("/signup", signupHdlrFunc(s))
	r.POST("/login", loginHdlrFunc(s))
	r.POST("/logout", logoutHldrFunc(s))
}

func signupHdlrFunc(s AuthServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrMsgInvalidRequestBody})
			return
		}

		var (
			tokens any
			userID string
		)
		err := s.WithTx(c.Request.Context(), func(ctx context.Context, tx store.Tx) error {
			var err error
			tokens, userID, err = s.AuthService().Signup(ctx, tx, req.Email, req.Password)
			if err != nil {
				return err
			}
			if errBot := s.BotService().CreateBot(ctx, tx, req.BotName, userID); errBot != nil {
				return fmt.Errorf("httphdlr.CreateBot: %w", errBot)
			}
			return nil
		})
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			switch {
			case errors.Is(err, authsvc.ErrUserExists):
				c.JSON(http.StatusConflict, gin.H{"error": authsvc.ErrUserExists.Error()})
			case errors.Is(err, authsvc.ErrInvalidCredentials):
				c.JSON(http.StatusBadRequest, gin.H{"error": authsvc.ErrInvalidCredentials.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			}
			return
		}
		c.JSON(http.StatusCreated, tokens)
	}
}

func loginHdlrFunc(s AuthServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrMsgInvalidRequestBody})
			return
		}
		tokens, err := s.AuthService().Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			switch err {
			case authsvc.ErrInvalidCredentials:
				c.JSON(http.StatusUnauthorized, gin.H{"error": authsvc.ErrInvalidCredentials.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
			}
			return
		}
		c.JSON(http.StatusOK, tokens)
	}
}

func logoutHldrFunc(s AuthServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req logoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrMsgInvalidRequestBody})
			return
		}
		if err := s.AuthService().Logout(c.Request.Context(), req.RefreshToken); err != nil {
			slog.Error(errutil.FormatErrChain(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": authsvc.ErrInvalidRefreshToken.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": authsvc.ErrLoggedOut.Error()})
	}
}
