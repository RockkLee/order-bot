package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/store"
)

type TxStore struct {
	db *sql.DB
}

func NewTxStore(db *sql.DB) *TxStore {
	if db == nil {
		panic("sqldb.NewTxStore(), the db ptr is nil")
	}
	return &TxStore{db: db}
}

func (s *TxStore) BeginTx(ctx context.Context) (store.Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("sqldb.TxStore.BeginTx: %w", err)
	}
	return tx, nil
}
