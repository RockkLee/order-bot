package menusvc

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/orderbotsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/resource"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
)

type Svc struct {
	menuStore          store.Menu
	menuItemStore      store.MenuItem
	publishedMenuStore *orderbotsqldb.PublishedMenuStore
	db                 *sqldb.DB
	orderBotDb         *sqldb.DB
	ctxFunc            util.CtxFunc
}

func NewSvc(
	rsrc *resource.Resource,
	ctxFunc util.CtxFunc,
	menuStore store.Menu,
	menuItemStore store.MenuItem,
	publishedMenuStore *orderbotsqldb.PublishedMenuStore,
) *Svc {
	if menuStore == nil || menuItemStore == nil || rsrc == nil || rsrc.DB == nil || rsrc.OrderBotDB == nil || publishedMenuStore == nil {
		panic("menusvc.NewSvc(), menuStore, menuItemStore, publishedMenuStore, db, or orderBotDb is nil")
	}
	return &Svc{
		menuStore:          menuStore,
		menuItemStore:      menuItemStore,
		publishedMenuStore: publishedMenuStore,
		db:                 rsrc.DB,
		orderBotDb:         rsrc.OrderBotDB,
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
			return fmt.Errorf("menusvc.GetMenuMenuItems(), duplicated bot ID: %w", ErrInvalidMenu)
		case !errors.Is(errFinding, store.ErrMenuNotFound):
			return fmt.Errorf("menusvc.GetMenuMenuItems: %w", errFinding)
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

func (s *Svc) GetMenu(ctx context.Context, botId string) (entities.Menu, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	menu, err := s.menuStore.FindByBotID(ctx, botId)
	if err != nil {
		return entities.Menu{}, fmt.Errorf("menusvc.GetMenu: %w", err)
	}
	return menu, nil
}

func (s *Svc) GetMenuMenuItems(ctx context.Context, botId string) (entities.Menu, []entities.MenuItem, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	menu, err := s.menuStore.FindByBotID(ctx, botId)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenuMenuItems: %w", err)
	}
	items, err := s.menuItemStore.FindItems(ctx, menu.ID)
	if err != nil {
		return entities.Menu{}, nil, fmt.Errorf("menusvc.GetMenuMenuItems: %w", err)
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
	menu, items, err := s.GetMenuMenuItems(ctx, botID)
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

func (s *Svc) IsMenuPublished(ctx context.Context, menuID string) (bool, error) {
	ctx, cancel := util.CallCtxFunc(ctx, s.ctxFunc)
	defer cancel()
	exists, err := s.publishedMenuStore.IsMenuPublished(ctx, menuID)
	if err != nil {
		return false, fmt.Errorf("menusvc.IsMenuPublished: %w", err)
	}
	return exists, nil
}
