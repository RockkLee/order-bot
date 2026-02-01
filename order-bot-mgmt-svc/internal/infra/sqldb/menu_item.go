package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
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
	insertMenuItemQuery     = `INSERT INTO menu_item (id, menu_id, menu_item_name) VALUES ($1, $2, $3);`
	selectMenuItemsByMenuID = `SELECT id, menu_id, menu_item_name FROM menu_item WHERE menu_id = $1 ORDER BY id;`
	deleteMenuItemsByMenuID = `DELETE FROM menu_item WHERE menu_id = $1;`
)

func NewMenuItemStore(db *pqsqldb.DB) *MenuItemStore {
	if db == nil {
		panic("sqldb.NewMenuItem(), the db ptr is nil")
	}
	return &MenuItemStore{db: db.Conn()}
}

func (s *MenuItemStore) FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error) {
	rows, errQry := s.db.QueryContext(ctx, selectMenuItemsByMenuID, menuID)
	if errQry != nil {
		return nil, fmt.Errorf("sqldb.MenuItemStore.FindItems: %w", errQry)
	}
	defer rows.Close()
	var items []entities.MenuItem
	for rows.Next() {
		var record MenuItemRecord
		if err := rows.Scan(&record.ID, &record.MenuID, &record.MenuItemName); err != nil {
			return nil, fmt.Errorf("sqldb.MenuItemStore.FindItems(), Scan: %w", err)
		}
		items = append(items, record.ToModel())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqldb.MenuItemStore.FindItems: %w", err)
	}
	return items, nil
}

func (s *MenuItemStore) DeleteMenuItems(ctx context.Context, tx store.Tx, menuID string) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.DeleteMenuItems: %w", store.ErrInvalidTx)
	}
	_, err := sqlTx.ExecContext(ctx, deleteMenuItemsByMenuID, menuID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.DeleteMenuItems(), ExecContext: %w", err)
	}
	return nil
}

func (s *MenuItemStore) CreateMenuItems(
	ctx context.Context, tx store.Tx,
	items []entities.MenuItem,
) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.CreateMenuItems(): %w", store.ErrInvalidTx)
	}
	for _, item := range items {
		record := MenuItemRecordFromModel(item)
		if _, err := sqlTx.ExecContext(ctx, insertMenuItemQuery, record.ID, record.MenuID, record.MenuItemName); err != nil {
			return fmt.Errorf("sqldb.MenuItemStore.CreateMenuItems(), ExecContext: %w", err)
		}
	}
	return nil
}
