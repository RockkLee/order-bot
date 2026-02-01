package menusvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
)

type Svc struct {
	menuStore     store.Menu
	menuItemStore store.MenuItem
	db            *pqsqldb.DB
	ctxFunc       util.CtxFunc
}

func NewSvc(db *pqsqldb.DB, ctxFunc util.CtxFunc, menuStore store.Menu, menuItemStore store.MenuItem) *Svc {
	if menuStore == nil || menuItemStore == nil || db == nil {
		panic("menusvc.NewSvc(), menuStore, menuItemStore or db is nil")
	}
	return &Svc{
		menuStore:     menuStore,
		menuItemStore: menuItemStore,
		db:            db,
		ctxFunc:       ctxFunc,
	}
}

func (s *Svc) CreateMenu(ctx context.Context, botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
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
			return ErrInvalidMenu
		case !errors.Is(errFinding, sql.ErrNoRows):
			return fmt.Errorf("menusvc.GetMenu: %w", errFinding)
		}
		menu = entities.Menu{
			ID:    util.NewID(),
			BotID: botID,
		}
		if err := s.menuStore.CreateMenu(ctx, tx, menu); err != nil {
			return fmt.Errorf("menusvc.CreateMenu: %w", err)
		}
		items = buildMenuItems(menu.ID, itemNames)
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
	items, err := s.menuItemStore.FindItems(ctx, botId)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenu: %w", err)
	}
	return menu, items, nil
}

func (s *Svc) UpdateMenu(ctx context.Context, botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	var (
		menu  entities.Menu
		items []entities.MenuItem
	)
	err := s.db.WithTx(ctx, func(ctx context.Context, tx store.Tx) error {
		var err error
		menu, err = s.menuStore.FindByBotID(ctx, botID)
		if err != nil {
			return fmt.Errorf("menusvc.UpdateMenu: %w", err)
		}
		items = buildMenuItems(menu.ID, itemNames)
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

func buildMenuItems(menuID string, names []string) []entities.MenuItem {
	items := make([]entities.MenuItem, 0, len(names))
	for _, name := range names {
		if name == "" {
			continue
		}
		items = append(items, entities.MenuItem{
			ID:           util.NewID(),
			MenuID:       menuID,
			MenuItemName: name,
		})
	}
	return items
}
