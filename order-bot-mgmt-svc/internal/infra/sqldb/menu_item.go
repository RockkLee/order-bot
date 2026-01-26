package sqldb

import (
	"context"
	"database/sql"
	"order-bot-mgmt-svc/internal/models/entities"
)

type MenuItemRecord struct {
	ID           string
	MenuID       string
	MenuItemName string
}

func MenuItemRecordFromModel(item entities.MenuItem) MenuItemRecord {
	return MenuItemRecord{
		ID:           item.ID,
		MenuID:       item.MenuID,
		MenuItemName: item.MenuItemName,
	}
}

func (r MenuItemRecord) ToModel() entities.MenuItem {
	return entities.MenuItem{
		ID:           r.ID,
		MenuID:       r.MenuID,
		MenuItemName: r.MenuItemName,
	}
}

type MenuItemStore struct {
	db *sql.DB
}

const (
	insertMenuItemQueryStandalone     = `INSERT INTO menu_item (id, menu_id, menu_item_name) VALUES ($1, $2, $3);`
	selectMenuItemsByMenuIDStandalone = `SELECT id, menu_id, menu_item_name FROM menu_item WHERE menu_id = $1 ORDER BY id;`
	deleteMenuItemsByMenuIDStandalone = `DELETE FROM menu_item WHERE menu_id = $1;`
)

func NewMenuItemStore(db *sql.DB) *MenuItemStore {
	return &MenuItemStore{db: db}
}

func (s *MenuItemStore) Create(ctx context.Context, item entities.MenuItem) error {
	_, err := s.db.ExecContext(ctx, insertMenuItemQueryStandalone, item.ID, item.MenuID, item.MenuItemName)
	return err
}

func (s *MenuItemStore) FindByMenuID(ctx context.Context, menuID string) ([]entities.MenuItem, error) {
	rows, err := s.db.QueryContext(ctx, selectMenuItemsByMenuIDStandalone, menuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entities.MenuItem
	for rows.Next() {
		var item entities.MenuItem
		if err := rows.Scan(&item.ID, &item.MenuID, &item.MenuItemName); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *MenuItemStore) DeleteByMenuID(ctx context.Context, menuID string) error {
	_, err := s.db.ExecContext(ctx, deleteMenuItemsByMenuIDStandalone, menuID)
	return err
}
