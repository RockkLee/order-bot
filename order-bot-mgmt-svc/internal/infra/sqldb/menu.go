package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type MenuRecord struct {
	ID    string
	BotID string
}

func MenuRecordFromModel(menu entities.Menu) MenuRecord {
	return MenuRecord{
		ID:    menu.ID,
		BotID: menu.BotID,
	}
}

func (r MenuRecord) ToModel() entities.Menu {
	return entities.Menu{
		ID:    r.ID,
		BotID: r.BotID,
	}
}

type MenuStore struct {
	db *sql.DB
}

const (
	insertMenuQuery         = `INSERT INTO menu (id, bot_id) VALUES ($1, $2);`
	selectMenuByIDQuery     = `SELECT id, bot_id FROM menu WHERE id = $1;`
	updateMenuQuery         = `UPDATE menu SET bot_id = $2 WHERE id = $1;`
	deleteMenuQuery         = `DELETE FROM menu WHERE id = $1;`
	insertMenuItemQuery     = `INSERT INTO menu_item (id, menu_id, menu_item_name) VALUES ($1, $2, $3);`
	selectMenuItemsByMenuID = `SELECT id, menu_id, menu_item_name FROM menu_item WHERE menu_id = $1 ORDER BY id;`
	deleteMenuItemsByMenuID = `DELETE FROM menu_item WHERE menu_id = $1;`
)

func NewMenuStore(db *pqsqldb.DB) *MenuStore {
	if db == nil {
		panic("sqldb.NewMenuStore(), the db ptr is nil")
	}
	return &MenuStore{db: db.Conn()}
}

func (s *MenuStore) FindByID(ctx context.Context, menuID string) (entities.Menu, error) {
	var record MenuRecord
	err := s.db.QueryRowContext(ctx, selectMenuByIDQuery, menuID).Scan(&record.ID, &record.BotID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Menu{}, fmt.Errorf("sqldb.MenuStore.FindByID: %w", store.ErrMenuNotFound)
		}
		return entities.Menu{}, fmt.Errorf("sqldb.MenuStore.FindByID: %w", err)
	}
	return record.ToModel(), nil
}

func (s *MenuStore) FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error) {
	rows, err := s.db.QueryContext(ctx, selectMenuItemsByMenuID, menuID)
	if err != nil {
		return nil, fmt.Errorf("sqldb.MenuStore.FindItems: %w", err)
	}
	defer rows.Close()
	var items []entities.MenuItem
	for rows.Next() {
		var record MenuItemRecord
		if err := rows.Scan(&record.ID, &record.MenuID, &record.MenuItemName); err != nil {
			return nil, fmt.Errorf("sqldb.MenuStore.FindItems: %w", err)
		}
		items = append(items, record.ToModel())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqldb.MenuStore.FindItems: %w", err)
	}
	return items, nil
}

func (s *MenuStore) CreateMenu(ctx context.Context, tx store.Tx, menu entities.Menu) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.CreateMenu(): %w", store.ErrInvalidTx)
	}
	record := MenuRecordFromModel(menu)
	_, err := sqlTx.ExecContext(ctx, insertMenuQuery, record.ID, record.BotID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.CreateMenu: %w", err)
	}
	return nil
}

func (s *MenuStore) UpdateMenu(ctx context.Context, tx store.Tx, menu entities.Menu) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.UpdateMenu(): %w", store.ErrInvalidTx)
	}
	record := MenuRecordFromModel(menu)
	result, err := sqlTx.ExecContext(ctx, updateMenuQuery, record.ID, record.BotID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu: %w", store.ErrMenuNotFound)
	}
	return nil
}

func (s *MenuStore) DeleteMenu(ctx context.Context, tx store.Tx, menuID string) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.DeleteMenu(): %w", store.ErrInvalidTx)
	}
	result, err := sqlTx.ExecContext(ctx, deleteMenuQuery, menuID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", store.ErrMenuNotFound)
	}
	return nil
}

func (s *MenuStore) DeleteMenuItems(ctx context.Context, tx store.Tx, menuID string) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.DeleteMenuItems(): %w", store.ErrInvalidTx)
	}
	_, err := sqlTx.ExecContext(ctx, deleteMenuItemsByMenuID, menuID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenuItems: %w", err)
	}
	return nil
}

func (s *MenuStore) CreateMenuItems(
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
			return fmt.Errorf("sqldb.MenuStore.CreateMenuItems: %w", err)
		}
	}
	return nil
}
