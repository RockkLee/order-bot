package sqldb

import (
	"context"
	"database/sql"
	"errors"
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
	insertMenuQuery         = `INSERT INTO menus (id, bot_id) VALUES ($1, $2);`
	selectMenuByIDQuery     = `SELECT id, bot_id FROM menus WHERE id = $1;`
	updateMenuQuery         = `UPDATE menus SET bot_id = $2 WHERE id = $1;`
	deleteMenuQuery         = `DELETE FROM menus WHERE id = $1;`
	insertMenuItemQuery     = `INSERT INTO menu_items (id, menu_id, menu_item_name) VALUES ($1, $2, $3);`
	selectMenuItemsByMenuID = `SELECT id, menu_id, menu_item_name FROM menu_items WHERE menu_id = $1 ORDER BY id;`
	deleteMenuItemsByMenuID = `DELETE FROM menu_items WHERE menu_id = $1;`
)

func NewMenuStore(db *sql.DB) *MenuStore {
	return &MenuStore{db: db}
}

func (s *MenuStore) BeginTx(ctx context.Context) (store.MenuTx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &MenuStore{db: s.db, tx: tx}, nil
}

func (s *MenuStore) FindByID(ctx context.Context, menuID string) (entities.Menu, error) {
	var record MenuRecord
	err := s.db.QueryRowContext(ctx, selectMenuByIDQuery, menuID).Scan(&record.ID, &record.BotID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Menu{}, store.ErrMenuNotFound
		}
		return entities.Menu{}, err
	}
	return record.ToModel(), nil
}

func (s *MenuStore) FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error) {
	rows, err := s.db.QueryContext(ctx, selectMenuItemsByMenuID, menuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entities.MenuItem
	for rows.Next() {
		var record MenuItemRecord
		if err := rows.Scan(&record.ID, &record.MenuID, &record.MenuItemName); err != nil {
			return nil, err
		}
		items = append(items, record.ToModel())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *MenuStore) CreateMenu(ctx context.Context, menu entities.Menu) error {
	if s.tx == nil {
		return errors.New("menu transaction not initialized")
	}
	record := MenuRecordFromModel(menu)
	_, err := s.tx.ExecContext(ctx, insertMenuQuery, record.ID, record.BotID)
	return err
}

func (s *MenuStore) UpdateMenu(ctx context.Context, menu entities.Menu) error {
	if s.tx == nil {
		return errors.New("menu transaction not initialized")
	}
	record := MenuRecordFromModel(menu)
	result, err := s.tx.ExecContext(ctx, updateMenuQuery, record.ID, record.BotID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return store.ErrMenuNotFound
	}
	return nil
}

func (s *MenuStore) DeleteMenu(ctx context.Context, menuID string) error {
	if s.tx == nil {
		return errors.New("menu transaction not initialized")
	}
	result, err := s.tx.ExecContext(ctx, deleteMenuQuery, menuID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return store.ErrMenuNotFound
	}
	return nil
}

func (s *MenuStore) DeleteMenuItems(ctx context.Context, menuID string) error {
	if s.tx == nil {
		return errors.New("menu transaction not initialized")
	}
	_, err := s.tx.ExecContext(ctx, deleteMenuItemsByMenuID, menuID)
	return err
}

func (s *MenuStore) CreateMenuItems(ctx context.Context, items []entities.MenuItem) error {
	if s.tx == nil {
		return errors.New("menu transaction not initialized")
	}
	for _, item := range items {
		record := MenuItemRecordFromModel(item)
		if _, err := s.tx.ExecContext(ctx, insertMenuItemQuery, record.ID, record.MenuID, record.MenuItemName); err != nil {
			return err
		}
	}
	return nil
}

func (s *MenuStore) Commit() error {
	if s.tx == nil {
		return errors.New("menu transaction not initialized")
	}
	return s.tx.Commit()
}

func (s *MenuStore) Rollback() error {
	if s.tx == nil {
		return errors.New("menu transaction not initialized")
	}
	return s.tx.Rollback()
}
