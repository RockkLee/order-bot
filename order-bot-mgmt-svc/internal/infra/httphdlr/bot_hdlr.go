package httphdlr

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/errutil"
	"order-bot-mgmt-svc/internal/util/jwtutil"

	"github.com/gin-gonic/gin"
)

type BotServer interface {
	BotService() *botsvc.Svc
	MenuService() *menusvc.Svc
	GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error)
}

const BotPrefix = "/bot"

func RegisterBotRoutes(r gin.IRoutes, s BotServer) {
	r.GET("/", getBotHdlrFunc(s))
}

func getBotHdlrFunc(s BotServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := jwtutil.GetTokenGin(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}
		// botIDAny, err := s.GetWithTx(c.Request.Context(), func(ctx context.Context, tx store.Tx) (any, error) {
		// 	return s.BotService().GetBotId(ctx, tx, token)
		// })
		// if err != nil {
		// 	slog.Error(errutil.FormatErrChain(err))
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load bot"})
		// 	return
		// }
		// botID, ok := botIDAny.(string)
		botID, err := s.BotService().GetBotId(nil, token)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load bot"})
			return
		}
		menu, err := s.MenuService().GetMenu(nil, botID)
		if err != nil {
			if errors.Is(err, store.ErrMenuNotFound) {
				c.JSON(http.StatusOK, gin.H{"bot_id": botID, "menu_id": nil})
				return
			}
			slog.Error(errutil.FormatErrChain(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load bot"})
			return
		}
		if !ok {
			slog.Error("bot id has unexpected type")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load bot"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"bot_id": botID, "menu_id": menu.ID})
	}
}
