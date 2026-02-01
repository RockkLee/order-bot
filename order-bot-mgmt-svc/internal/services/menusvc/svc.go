package menusvc

import (
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

func (s *Svc) CreateMenu(botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := s.ctxFunc()
	defer cancel()
	tx, errTx := s.db.BeginTx(ctx)
	if errTx != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.CreateMenu: %w", errTx)
	}
	if _, err := s.menuStore.FindByBotID(ctx, botID); !errors.Is(err, sql.ErrNoRows) {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenu: %w", err)
	}
	menu := entities.Menu{
		ID:    util.NewID(),
		BotID: botID,
	}
	if err := s.menuStore.CreateMenu(ctx, tx, menu); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, fmt.Errorf("menusvc.CreateMenu: %w", err)
	}
	items := buildMenuItems(menu.ID, itemNames)
	if err := s.menuItemStore.CreateMenuItems(ctx, tx, items); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, fmt.Errorf("menusvc.CreateMenu: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.CreateMenu(), Commit: %w", err)
	}
	return menu, items, nil
}

func (s *Svc) GetMenu(botId string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := s.ctxFunc()
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

func (s *Svc) UpdateMenu(botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := s.ctxFunc()
	defer cancel()
	tx, errTx := s.db.BeginTx(ctx)
	if errTx != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", errTx)
	}
	menu, errMenu := s.menuStore.FindByBotID(ctx, botID)
	if errMenu != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", errMenu)
	}
	items := buildMenuItems(menu.ID, itemNames)
	if err := s.menuItemStore.CreateMenuItems(ctx, tx, items); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", err)
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
