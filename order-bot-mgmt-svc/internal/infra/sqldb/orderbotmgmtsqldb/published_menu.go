package orderbotmgmtsqldb

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/sqldbexecutor"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type PublishedMenuStore struct {
	db *sql.DB
}

const (
	deletePublishedMenuItemsQuery = `DELETE FROM menu_items;`
	insertPublishedMenuItemQuery  = `INSERT INTO menu_items (id, sku, name, description, price_cents, is_available, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);`
)

func NewPublishedMenuStore(db *pqsqldb.DB) *PublishedMenuStore {
	if db == nil {
		panic("orderbotmgmtsqldb.NewPublishedMenuStore(), the db ptr is nil")
	}
	return &PublishedMenuStore{db: db.Conn()}
}

func (s *PublishedMenuStore) ReplaceMenuItems(ctx context.Context, tx store.Tx, items []entities.MenuItem) error {
	execer, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems: %w", err)
	}
	if _, err := execer.ExecContext(ctx, deletePublishedMenuItemsQuery); err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), delete: %w", err)
	}
	for _, item := range items {
		record := PublishedMenuItemRecordFromModel(item)
		if _, err := execer.ExecContext(
			ctx,
			insertPublishedMenuItemQuery,
			record.ID,
			record.SKU,
			record.Name,
			record.Description,
			record.PriceCents,
			record.IsAvailable,
			record.CreatedAt,
		); err != nil {
			return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), insert: %w", err)
		}
	}
	return nil
}
