package menusvc

import (
	"context"
	"errors"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"time"
)

const menuQueryTimeout = 2 * time.Second

type Svc struct {
	menuStore store.Menu
}

func NewSvc(menuStore store.Menu) *Svc {
	return &Svc{menuStore: menuStore}
}

func (s *Svc) CreateMenu(botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
	if s.menuStore == nil {
		return entities.Menu{}, nil, errors.New("menu store not configured")
	}
	if botID == "" {
		return entities.Menu{}, nil, ErrInvalidMenu
	}
	ctx, cancel := s.menuContext()
	defer cancel()
	tx, err := s.menuStore.BeginTx(ctx)
	if err != nil {
		return entities.Menu{}, nil, err
	}
	menu := entities.Menu{
		ID:    util.NewID(),
		BotID: botID,
	}
	if err := tx.CreateMenu(ctx, menu); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, err
	}
	items := buildMenuItems(menu.ID, itemNames)
	if err := tx.CreateMenuItems(ctx, items); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Menu{}, nil, err
	}
	return menu, items, nil
}

func (s *Svc) GetMenu(menuID string) (entities.Menu, []entities.MenuItem, error) {
	if s.menuStore == nil {
		return entities.Menu{}, nil, errors.New("menu store not configured")
	}
	if menuID == "" {
		return entities.Menu{}, nil, ErrInvalidMenu
	}
	ctx, cancel := s.menuContext()
	defer cancel()
	menu, err := s.menuStore.FindByID(ctx, menuID)
	if err != nil {
		return entities.Menu{}, nil, err
	}
	items, err := s.menuStore.FindItems(ctx, menuID)
	if err != nil {
		return entities.Menu{}, nil, err
	}
	return menu, items, nil
}

func (s *Svc) UpdateMenu(menuID, botID string, itemNames []string) (entities.Menu, []entities.MenuItem, error) {
	if s.menuStore == nil {
		return entities.Menu{}, nil, errors.New("menu store not configured")
	}
	if menuID == "" || botID == "" {
		return entities.Menu{}, nil, ErrInvalidMenu
	}
	ctx, cancel := s.menuContext()
	defer cancel()
	tx, err := s.menuStore.BeginTx(ctx)
	if err != nil {
		return entities.Menu{}, nil, err
	}
	menu := entities.Menu{
		ID:    menuID,
		BotID: botID,
	}
	if err := tx.UpdateMenu(ctx, menu); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, err
	}
	if err := tx.DeleteMenuItems(ctx, menuID); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, err
	}
	items := buildMenuItems(menuID, itemNames)
	if err := tx.CreateMenuItems(ctx, items); err != nil {
		_ = tx.Rollback()
		return entities.Menu{}, nil, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Menu{}, nil, err
	}
	return menu, items, nil
}

func (s *Svc) DeleteMenu(menuID string) error {
	if s.menuStore == nil {
		return errors.New("menu store not configured")
	}
	if menuID == "" {
		return ErrInvalidMenu
	}
	ctx, cancel := s.menuContext()
	defer cancel()
	tx, err := s.menuStore.BeginTx(ctx)
	if err != nil {
		return err
	}
	if err := tx.DeleteMenuItems(ctx, menuID); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.DeleteMenu(ctx, menuID); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
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

func (s *Svc) menuContext() (context.Context, context.CancelFunc) {
	return services.QueryContext(menuQueryTimeout)
}
