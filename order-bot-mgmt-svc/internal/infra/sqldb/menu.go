package sqldb

import (
	"context"
	"database/sql"
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
	insertMenuQuery        = `INSERT INTO menu (id, bot_id) VALUES ($1, $2);`
	selectMenuByBotIDQuery = `SELECT id, bot_id FROM menu WHERE id = $1;`
	updateMenuQuery        = `UPDATE menu SET bot_id = $2 WHERE id = $1;`
	deleteMenuQuery        = `DELETE FROM menu WHERE id = $1;`
)

func NewMenuStore(db *pqsqldb.DB) *MenuStore {
	if db == nil {
		panic("sqldb.NewMenuStore(), the db ptr is nil")
	}
	return &MenuStore{db: db.Conn()}
}

func (s *MenuStore) FindByBotID(ctx context.Context, menuID string) (entities.Menu, error) {
	var record MenuRecord
	err := s.db.QueryRowContext(ctx, selectMenuByBotIDQuery, menuID).Scan(&record.ID, &record.BotID)
	if err != nil {
		return entities.Menu{}, fmt.Errorf("sqldb.MenuStore.FindByBotID: %w", err)
	}
	return record.ToModel(), nil
}

func (s *MenuStore) CreateMenu(ctx context.Context, tx store.Tx, menu entities.Menu) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.CreateMenu: %w", store.ErrInvalidTx)
	}
	record := MenuRecordFromModel(menu)
	_, err := sqlTx.ExecContext(ctx, insertMenuQuery, record.ID, record.BotID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.CreateMenu(), ExecContext: %w", err)
	}
	return nil
}

func (s *MenuStore) UpdateMenu(ctx context.Context, tx store.Tx, menu entities.Menu) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("sqldb.UpdateMenu: %w", store.ErrInvalidTx)
	}
	record := MenuRecordFromModel(menu)
	result, err := sqlTx.ExecContext(ctx, updateMenuQuery, record.ID, record.BotID)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu(), ExecContext: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu(), RowsAffected: %w", err)
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
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu(), ExecContext: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu(), RowsAffected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", store.ErrMenuNotFound)
	}
	return nil
}
