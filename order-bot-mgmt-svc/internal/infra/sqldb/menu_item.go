package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/sqldbexecutor"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type MenuItemRecord struct {
	ID           string
	MenuID       string
	MenuItemName string
	Price        float64
}

func MenuItemRecordFromModel(item entities.MenuItem) MenuItemRecord {
	return MenuItemRecord{
		ID:           item.ID,
		MenuID:       item.MenuID,
		MenuItemName: item.MenuItemName,
		Price:        item.Price,
	}
}

func (r MenuItemRecord) ToModel() entities.MenuItem {
	return entities.MenuItem{
		ID:           r.ID,
		MenuID:       r.MenuID,
		MenuItemName: r.MenuItemName,
		Price:        r.Price,
	}
}

type MenuItemStore struct {
	db *sql.DB
}

const (
	insertMenuItemQuery     = `INSERT INTO menu_item (id, menu_id, menu_item_name, price) VALUES ($1, $2, $3, $4);`
	selectMenuItemsByMenuID = `SELECT id, menu_id, menu_item_name, price FROM menu_item WHERE menu_id = $1 ORDER BY id;`
	deleteMenuItemsByMenuID = `DELETE FROM menu_item WHERE menu_id = $1;`
)

func NewMenuItemStore(db *pqsqldb.DB) *MenuItemStore {
	if db == nil {
		panic("sqldb.NewMenuItem(), the db ptr is nil")
	}
	return &MenuItemStore{db: db.Conn()}
}

func (s *MenuItemStore) FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error) {
	execer, err := sqldbexecutor.Executor(s.db, nil)
	if err != nil {
		return nil, fmt.Errorf("sqldb.MenuItemStore.FindItems: %w", err)
	}
	rows, err := execer.QueryContext(ctx, selectMenuItemsByMenuID, menuID)
	if err != nil {
		return nil, fmt.Errorf("sqldb.MenuItemStore.FindItems: %w", err)
	}
	defer rows.Close()
	var items []entities.MenuItem
	for rows.Next() {
		var record MenuItemRecord
		if err := rows.Scan(&record.ID, &record.MenuID, &record.MenuItemName, &record.Price); err != nil {
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
	execer, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.DeleteMenuItems: %w", err)
	}
	_, err = execer.ExecContext(ctx, deleteMenuItemsByMenuID, menuID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.DeleteMenuItems(), ExecContext: %w", err)
	}
	return nil
}

func (s *MenuItemStore) CreateMenuItems(
	ctx context.Context, tx store.Tx,
	items []entities.MenuItem,
) error {
	execer, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.CreateMenuItems: %w", err)
	}
	for _, item := range items {
		record := MenuItemRecordFromModel(item)
		if _, err := execer.ExecContext(ctx, insertMenuItemQuery, record.ID, record.MenuID, record.MenuItemName, record.Price); err != nil {
			return fmt.Errorf("sqldb.MenuItemStore.CreateMenuItems(), ExecContext: %w", err)
		}
	}
	return nil
}
