package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/store"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore(db *sql.DB) *DBStore {
	if db == nil {
		panic("sqldb.NewDBStore(), the db ptr is nil")
	}
	return &DBStore{db: db}
}

func (s *DBStore) BeginTx(ctx context.Context) (store.Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("sqldb.DBStore.BeginTx: %w", err)
	}
	return tx, nil
}
