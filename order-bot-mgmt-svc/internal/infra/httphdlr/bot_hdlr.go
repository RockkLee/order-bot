package httphdlr

import (
	"context"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/errutil"
	"strings"

	"github.com/gin-gonic/gin"
)

type BotServer interface {
	BotService() *botsvc.Svc
	GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error)
}

const BotPrefix = "/bot"

func RegisterBotRoutes(r gin.IRoutes, s BotServer) {
	r.GET("/", getBotHdlrFunc(s))
}

func getBotHdlrFunc(s BotServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := tokenFromRequest(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}
		botIDAny, err := s.GetWithTx(c.Request.Context(), func(ctx context.Context, tx store.Tx) (any, error) {
			return s.BotService().GetBotId(ctx, tx, token)
		})
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load bot"})
			return
		}
		botID, ok := botIDAny.(string)
		if !ok {
			slog.Error("bot id has unexpected type")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load bot"})
			return
		}
		c.JSON(http.StatusOK, botID)
	}
}

func tokenFromRequest(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", false
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", false
	}
	accessToken := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	return accessToken, accessToken != ""
}
