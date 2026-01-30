package menusvc

import (
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"time"
)

const menuQueryTimeout = 2 * time.Second

type Svc struct {
	menuStore     store.Menu
	menuItemStore store.MenuItem
	db            *pqsqldb.DB
	ctxFunc       util.CtxFunc
}

func NewSvc(menuStore store.Menu, menuItemStore store.MenuItem, db *pqsqldb.DB) *Svc {
	if menuStore == nil || menuItemStore == nil || db == nil {
		panic("menusvc.NewSvc(), menuStore, menuItemStore or db is nil")
	}
	return &Svc{
		menuStore:     menuStore,
		menuItemStore: menuItemStore,
		db:            db,
		ctxFunc:       util.NewCtxFunc(menuQueryTimeout),
	}
}

func (s *Svc) CreateMenu(botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := s.ctxFunc()
	defer cancel()
	tx, errTx := s.db.BeginTx(ctx)
	if errTx != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.CreateMenu: %w", errTx)
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
		return entities.Menu{}, nil, fmt.Errorf("menusvc.CreateMenu: %w", err)
	}
	return menu, items, nil
}

func (s *Svc) GetMenu(menuID string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := s.ctxFunc()
	defer cancel()
	menu, err := s.menuStore.FindByID(ctx, menuID)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenu: %w", err)
	}
	items, err := s.menuItemStore.FindItems(ctx, menuID)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenu: %w", err)
	}
	return menu, items, nil
}

func (s *Svc) UpdateMenu(menuID, botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := s.ctxFunc()
	defer cancel()
	tx, errTx := s.db.BeginTx(ctx)
	if errTx != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", errTx)
	}
	menu := entities.Menu{
		ID:    menuID,
		BotID: botID,
	}
	if err := s.menuStore.UpdateMenu(ctx, tx, menu); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", err)
	}
	if err := s.menuItemStore.DeleteMenuItems(ctx, tx, menuID); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", err)
	}
	items := buildMenuItems(menuID, itemNames)
	if err := s.menuItemStore.CreateMenuItems(ctx, tx, items); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.UpdateMenu: %w", err)
	}
	return menu, items, nil
}

func (s *Svc) DeleteMenu(menuID string) error {
	ctx, cancel := s.ctxFunc()
	defer cancel()
	tx, errTx := s.db.BeginTx(ctx)
	if errTx != nil {
		return fmt.Errorf("menusvc.DeleteMenu: %w", errTx)
	}
	if err := s.menuItemStore.DeleteMenuItems(ctx, tx, menuID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("menusvc.DeleteMenu: %w", err)
	}
	if err := s.menuStore.DeleteMenu(ctx, tx, menuID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("menusvc.DeleteMenu: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("menusvc.DeleteMenu: %w", err)
	}
	return nil
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
