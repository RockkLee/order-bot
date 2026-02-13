package botsvc

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"order-bot-mgmt-svc/internal/util/jwtutil"
)

type Svc struct {
	db           *sqldb.DB
	ctxFunc      util.CtxFunc
	botStore     store.Bot
	userBotStore store.UserBot
	accessSecret []byte
}

func NewSvc(db *sqldb.DB, ctxFunc util.CtxFunc, cfg config.Config, botStore store.Bot, userBotStore store.UserBot) *Svc {
	if botStore == nil || db == nil {
		panic("botsvc.NewSvc(), botStore, menuItemStore or db is nil")
	}
	return &Svc{
		botStore:     botStore,
		userBotStore: userBotStore,
		db:           db,
		ctxFunc:      ctxFunc,
		accessSecret: []byte(cfg.Auth.AccessSecret),
	}
}

func (s *Svc) CreateBot(ctx context.Context, tx store.Tx, name string, userId string) error {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	newBot := entities.Bot{
		ID:      util.NewID(),
		BotName: name,
	}
	if err := s.botStore.Create(ctx, tx, newBot); err != nil {
		return fmt.Errorf("botsvc.CreateBot: %w", err)
	}
	newUserBot := entities.UserBot{
		ID:     util.NewID(),
		UserID: userId,
		BotID:  newBot.ID,
	}
	if err := s.userBotStore.Create(ctx, tx, newUserBot); err != nil {
		return fmt.Errorf("botsvc.CreateBot: %w", err)
	}
	return nil
}

func (s *Svc) GetBotId(ctx context.Context, tx store.Tx, tokenStr string) (botId string, err error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()

	claims, err := jwtutil.ParseJWT(s.accessSecret, tokenStr)
	if err != nil {
		return "", fmt.Errorf("botsvc.GetBotId(): %w", err)
	}
	userBots, err := s.userBotStore.FindByUserID(ctx, tx, claims.Sub)
	if err != nil {
		return "", fmt.Errorf("botsvc.GetBotId: %w", err)
	}
	return userBots[0].BotID, err
}
