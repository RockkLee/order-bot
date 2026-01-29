package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	tx *sql.Tx
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

func NewMenuStore(db *sql.DB) *MenuStore {
	if db == nil {
		panic("sqldb.NewMenuStore(), the db ptr is nil")
	}
	return &MenuStore{db: db}
}

func (s *MenuStore) BeginTx(ctx context.Context) (store.MenuTx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("sqldb.MenuStore.BeginTx: %w", err)
	}
	return &MenuStore{db: s.db, tx: tx}, nil
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

func (s *MenuStore) CreateMenu(ctx context.Context, menu entities.Menu) error {
	if s.tx == nil {
		return fmt.Errorf("sqldb.MenuStore.CreateMenu: %w", errors.New("menu transaction not initialized"))
	}
	record := MenuRecordFromModel(menu)
	_, err := s.tx.ExecContext(ctx, insertMenuQuery, record.ID, record.BotID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.CreateMenu: %w", err)
	}
	return nil
}

func (s *MenuStore) UpdateMenu(ctx context.Context, menu entities.Menu) error {
	if s.tx == nil {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu: %w", errors.New("menu transaction not initialized"))
	}
	record := MenuRecordFromModel(menu)
	result, err := s.tx.ExecContext(ctx, updateMenuQuery, record.ID, record.BotID)
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

func (s *MenuStore) DeleteMenu(ctx context.Context, menuID string) error {
	if s.tx == nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", errors.New("menu transaction not initialized"))
	}
	result, err := s.tx.ExecContext(ctx, deleteMenuQuery, menuID)
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

func (s *MenuStore) DeleteMenuItems(ctx context.Context, menuID string) error {
	if s.tx == nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenuItems: %w", errors.New("menu transaction not initialized"))
	}
	_, err := s.tx.ExecContext(ctx, deleteMenuItemsByMenuID, menuID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenuItems: %w", err)
	}
	return nil
}

func (s *MenuStore) CreateMenuItems(ctx context.Context, items []entities.MenuItem) error {
	if s.tx == nil {
		return fmt.Errorf("sqldb.MenuStore.CreateMenuItems: %w", errors.New("menu transaction not initialized"))
	}
	for _, item := range items {
		record := MenuItemRecordFromModel(item)
		if _, err := s.tx.ExecContext(ctx, insertMenuItemQuery, record.ID, record.MenuID, record.MenuItemName); err != nil {
			return fmt.Errorf("sqldb.MenuStore.CreateMenuItems: %w", err)
		}
	}
	return nil
}

func (s *MenuStore) Commit() error {
	if s.tx == nil {
		return fmt.Errorf("sqldb.MenuStore.Commit: %w", errors.New("menu transaction not initialized"))
	}
	if err := s.tx.Commit(); err != nil {
		return fmt.Errorf("sqldb.MenuStore.Commit: %w", err)
	}
	return nil
}

func (s *MenuStore) Rollback() error {
	if s.tx == nil {
		return fmt.Errorf("sqldb.MenuStore.Rollback: %w", errors.New("menu transaction not initialized"))
	}
	if err := s.tx.Rollback(); err != nil {
		return fmt.Errorf("sqldb.MenuStore.Rollback: %w", err)
	}
	return nil
}
