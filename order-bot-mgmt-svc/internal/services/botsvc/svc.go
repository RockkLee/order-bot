package botsvc

import (
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

func (s *Svc) CreateBot(name string, userId string) error {
	ctx, cancel := s.ctxFunc()
	defer cancel()
	newBot := entities.Bot{
		ID:      util.NewID(),
		BotName: name,
	}
	if err := s.botStore.Create(ctx, newBot); err != nil {
		return fmt.Errorf("bitsvc.CreateBot: %w", err)
	}
	newUserBot := entities.UserBot{
		ID:     util.NewID(),
		UserID: userId,
		BotID:  newBot.ID,
	}
	if err := s.userBotStore.Create(ctx, newUserBot); err != nil {
		return fmt.Errorf("bitsvc.CreateBot: %w", err)
	}
	return nil
}
