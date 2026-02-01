package menusvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/orderbotmgmtsqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
)

type Svc struct {
	menuStore          store.Menu
	menuItemStore      store.MenuItem
	publishedMenuStore *orderbotmgmtsqldb.PublishedMenuStore
	db                 *pqsqldb.DB
	orderBotDb         *pqsqldb.DB
	ctxFunc            util.CtxFunc
}

func NewSvc(
	db *pqsqldb.DB,
	orderBotDb *pqsqldb.DB,
	ctxFunc util.CtxFunc,
	menuStore store.Menu,
	menuItemStore store.MenuItem,
	publishedMenuStore *orderbotmgmtsqldb.PublishedMenuStore,
) *Svc {
	if menuStore == nil || menuItemStore == nil || db == nil || orderBotDb == nil || publishedMenuStore == nil {
		panic("menusvc.NewSvc(), menuStore, menuItemStore, publishedMenuStore, db, or orderBotDb is nil")
	}
	return &Svc{
		menuStore:          menuStore,
		menuItemStore:      menuItemStore,
		publishedMenuStore: publishedMenuStore,
		db:                 db,
		orderBotDb:         orderBotDb,
		ctxFunc:            ctxFunc,
	}
}

func (s *Svc) CreateMenu(ctx context.Context, botID string, menuItems []entities.MenuItem) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	var (
		menu  entities.Menu
		items []entities.MenuItem
	)
	err := s.db.WithTx(ctx, func(ctx context.Context, tx store.Tx) error {
		_, errFinding := s.menuStore.FindByBotID(ctx, botID)
		switch {
		case errFinding == nil:
			return fmt.Errorf("menusvc.GetMenu(), duplicated bot ID: %w", ErrInvalidMenu)
		case !errors.Is(errFinding, sql.ErrNoRows):
			return fmt.Errorf("menusvc.GetMenu: %w", errFinding)
		}
		menu = entities.Menu{
			ID:    util.NewID(),
			BotID: botID,
		}
		items = menuItems
		for idx := range items {
			items[idx].MenuID = menu.ID
		}
		if err := s.menuStore.CreateMenu(ctx, tx, menu); err != nil {
			return fmt.Errorf("menusvc.CreateMenu: %w", err)
		}
		if err := s.menuItemStore.CreateMenuItems(ctx, tx, items); err != nil {
			return fmt.Errorf("menusvc.CreateMenu: %w", err)
		}
		return nil
	})
	if err != nil {
		return entities.Menu{}, nil, err
	}
	return menu, items, nil
}

func (s *Svc) GetMenu(ctx context.Context, botId string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	menu, err := s.menuStore.FindByBotID(ctx, botId)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenu: %w", err)
	}
	items, err := s.menuItemStore.FindItems(ctx, menu.ID)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenu: %w", err)
	}
	return menu, items, nil
}

func (s *Svc) UpdateMenu(ctx context.Context, botID string, items []entities.MenuItem) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	var menu entities.Menu
	err := s.db.WithTx(ctx, func(ctx context.Context, tx store.Tx) error {
		menu, errMenu := s.menuStore.FindByBotID(ctx, botID)
		for idx := range items {
			items[idx].MenuID = menu.ID
		}
		if errMenu != nil {
			return fmt.Errorf("menusvc.UpdateMenu: %w", errMenu)
		}
		if err := s.menuItemStore.DeleteMenuItems(ctx, tx, menu.ID); err != nil {
			return fmt.Errorf("menusvc.UpdateMenu: %w", err)
		}
		if err := s.menuItemStore.CreateMenuItems(ctx, tx, items); err != nil {
			return fmt.Errorf("menusvc.UpdateMenu: %w", err)
		}
		return nil
	})
	if err != nil {
		return entities.Menu{}, nil, err
	}
	return menu, items, nil
}

func (s *Svc) PublishMenu(ctx context.Context, botID string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	menu, items, err := s.GetMenu(ctx, botID)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.PublishMenu: %w", err)
	}
	if err := s.orderBotDb.WithTx(ctx, func(ctx context.Context, tx store.Tx) error {
		if err := s.publishedMenuStore.ReplaceMenuItems(ctx, tx, menu, items); err != nil {
			return fmt.Errorf("menusvc.PublishMenu: %w", err)
		}
		return nil
	}); err != nil {
		return entities.Menu{}, nil, err
	}
	return menu, items, nil
}
