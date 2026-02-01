package botsvc

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
)

type Svc struct {
	botStore     store.Bot
	userBotStore store.UserBot
	db           *pqsqldb.DB
	ctxFunc      util.CtxFunc
}

func NewSvc(botStore store.Bot, userBotStore store.UserBot, db *pqsqldb.DB, ctxFunc util.CtxFunc) *Svc {
	if botStore == nil || db == nil {
		panic("menusvc.NewSvc(), botStore, menuItemStore or db is nil")
	}
	return &Svc{
		botStore:     botStore,
		userBotStore: userBotStore,
		db:           db,
		ctxFunc:      ctxFunc,
	}
}

func (s *Svc) CreateBot(ctx context.Context, tx store.Tx, name string, userId string) error {
	ctx, cancel := s.ctxWithFallback(ctx)
	defer cancel()
	newBot := entities.Bot{
		ID:      util.NewID(),
		BotName: name,
	}
	if err := s.botStore.Create(ctx, tx, newBot); err != nil {
		return fmt.Errorf("bitsvc.CreateBot: %w", err)
	}
	newUserBot := entities.UserBot{
		ID:     util.NewID(),
		UserID: userId,
		BotID:  newBot.ID,
	}
	if err := s.userBotStore.Create(ctx, tx, newUserBot); err != nil {
		return fmt.Errorf("bitsvc.CreateBot: %w", err)
	}
	return nil
}

func (s *Svc) ctxWithFallback(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return s.ctxFunc()
	}
	return ctx, func() {}
}
