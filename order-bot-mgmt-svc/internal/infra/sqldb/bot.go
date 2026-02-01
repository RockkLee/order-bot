package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb/pqsqldb"
	"order-bot-mgmt-svc/internal/infra/sqldb/sqldbexecutor"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type BotRecord struct {
	ID      string
	BotName string
}

func BotRecordFromModel(bot entities.Bot) BotRecord {
	return BotRecord{
		ID:      bot.ID,
		BotName: bot.BotName,
	}
}

func (r BotRecord) ToModel() entities.Bot {
	return entities.Bot{
		ID:      r.ID,
		BotName: r.BotName,
	}
}

type BotStore struct {
	db *sql.DB
}

const (
	insertBotQuery = `INSERT INTO bot (id, bot_name) VALUES ($1, $2);`
	selectBotByID  = `SELECT id, bot_name FROM bot WHERE id = $1;`
)

func NewBotStore(db *pqsqldb.DB) *BotStore {
	if db == nil {
		panic("sqldb.NewBotStore(), the db ptr is nil")
	}
	return &BotStore{db: db.Conn()}
}

func (s *BotStore) Create(ctx context.Context, tx store.Tx, bot entities.Bot) error {
	record := BotRecordFromModel(bot)
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.BotStore.Create(): %w", err)
	}
	_, err = exec.ExecContext(ctx, insertBotQuery, record.ID, record.BotName)
	if err != nil {
		return fmt.Errorf("sqldb.BotStore.Create(), ExecContext: %w", err)
	}
	return nil
}

func (s *BotStore) FindByID(ctx context.Context, tx store.Tx, id string) (entities.Bot, error) {
	var record BotRecord
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return entities.Bot{}, fmt.Errorf("sqldb.BotStore.FindByBotID: %w", err)
	}
	err = exec.QueryRowContext(ctx, selectBotByID, id).Scan(&record.ID, &record.BotName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Bot{}, fmt.Errorf("sqldb.BotStore.FindByBotID: %w", store.ErrBotNotFound)
		}
		return entities.Bot{}, fmt.Errorf("sqldb.BotStore.FindByBotID: %w", err)
	}
	return record.ToModel(), nil
}
