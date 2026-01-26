package sqldb

import (
	"context"
	"database/sql"
	"errors"
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
	insertBotQuery = `INSERT INTO bots (id, bot_name) VALUES ($1, $2);`
	selectBotByID  = `SELECT id, bot_name FROM bots WHERE id = $1;`
)

func NewBotStore(db *sql.DB) *BotStore {
	return &BotStore{db: db}
}

func (s *BotStore) Create(ctx context.Context, bot entities.Bot) error {
	record := BotRecordFromModel(bot)
	_, err := s.db.ExecContext(ctx, insertBotQuery, record.ID, record.BotName)
	return err
}

func (s *BotStore) FindByID(ctx context.Context, id string) (entities.Bot, error) {
	var record BotRecord
	err := s.db.QueryRowContext(ctx, selectBotByID, id).Scan(&record.ID, &record.BotName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Bot{}, store.ErrBotNotFound
		}
		return entities.Bot{}, err
	}
	return record.ToModel(), nil
}
