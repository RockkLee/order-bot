package httphdlrsold

import (
	"context"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/errutil"
	"order-bot-mgmt-svc/internal/util/jwtutil"
)

type BotServer interface {
	BotService() *botsvc.Svc
	GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error)
}

const BotPrefix = "/bot"

func BotHdlr(s BotServer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getBotHdlrFunc(s))
	return mux
}

func getBotHdlrFunc(s BotServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, done := jwtutil.GetToken(w, r)
		if done {
			return
		}
		botIdAny, err := s.GetWithTx(r.Context(), func(ctx context.Context, tx store.Tx) (any, error) {
			botId, err := s.BotService().GetBotId(ctx, tokenStr)
			if err != nil {
				return nil, err
			}
			return botId, nil
		})
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}
		botId, ok := botIdAny.(string)
		if !ok {
			slog.Error("bot id has unexpected type")
			http.Error(w, "failed to load bot", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, botId)
	}
}
