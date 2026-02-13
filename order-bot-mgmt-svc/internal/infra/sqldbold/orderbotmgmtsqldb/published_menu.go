package orderbotmgmtsqldb

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/infra/sqldbold/sqldbexecutor"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type PublishedMenuStore struct {
	db *sql.DB
}

const (
	deletePublishedMenuQuery      = `DELETE FROM published_menu WHERE bot_id = $1;`
	deletePublishedMenuItemsQuery = `DELETE FROM published_menu_item WHERE menu_id = $1;`
	insertPublishedMenuQuery      = `INSERT INTO published_menu (id, bot_id) VALUES ($1, $2);`
	insertPublishedMenuItemQuery  = `INSERT INTO published_menu_item (id, menu_id, menu_item_name, price) VALUES ($1, $2, $3, $4);`
)

func NewPublishedMenuStore(db *sqldb.DB) *PublishedMenuStore {
	if db == nil {
		panic("orderbotmgmtsqldb.NewPublishedMenuStore(), the db ptr is nil")
	}
	return &PublishedMenuStore{db: db.Conn()}
}

func (s *PublishedMenuStore) ReplaceMenuItems(
	ctx context.Context,
	tx store.Tx,
	menu entities.Menu,
	items []entities.MenuItem,
) error {
	execer, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems: %w", err)
	}
	if _, err := execer.ExecContext(ctx, deletePublishedMenuQuery, menu.BotID); err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), delete menu: %w", err)
	}
	if _, err := execer.ExecContext(ctx, insertPublishedMenuQuery, menu.ID, menu.BotID); err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), insert menu: %w", err)
	}
	if _, err := execer.ExecContext(ctx, deletePublishedMenuItemsQuery, menu.ID); err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), delete menu_item: %w", err)
	}
	for _, item := range items {
		if _, err := execer.ExecContext(
			ctx,
			insertPublishedMenuItemQuery,
			item.ID,
			item.MenuID,
			item.MenuItemName,
			item.Price,
		); err != nil {
			return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), insert: %w", err)
		}
	}
	return nil
}
