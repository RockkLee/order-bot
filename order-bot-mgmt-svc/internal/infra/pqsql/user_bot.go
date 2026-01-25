package pqsql

import (
	"context"
	"database/sql"
	"errors"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type UserBotRecord struct {
	ID     string
	UserID string
	BotID  string
}

func UserBotRecordFromModel(userBot entities.UserBot) UserBotRecord {
	return UserBotRecord{
		ID:     userBot.ID,
		UserID: userBot.UserID,
		BotID:  userBot.BotID,
	}
}

func (r UserBotRecord) ToModel() entities.UserBot {
	return entities.UserBot{
		ID:     r.ID,
		UserID: r.UserID,
		BotID:  r.BotID,
	}
}

type UserBotStore struct {
	db *sql.DB
}

const (
	insertUserBotQuery   = `INSERT INTO user_bots (id, user_id, bot_id) VALUES ($1, $2, $3);`
	selectUserBotsByUser = `SELECT id, user_id, bot_id FROM user_bots WHERE user_id = $1;`
)

func NewUserBotStore(db *sql.DB) *UserBotStore {
	return &UserBotStore{db: db}
}

func (s *UserBotStore) Create(ctx context.Context, userBot entities.UserBot) error {
	record := UserBotRecordFromModel(userBot)
	_, err := s.db.ExecContext(ctx, insertUserBotQuery, record.ID, record.UserID, record.BotID)
	return err
}

func (s *UserBotStore) FindByUserID(ctx context.Context, userID string) ([]entities.UserBot, error) {
	rows, err := s.db.QueryContext(ctx, selectUserBotsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []entities.UserBot
	for rows.Next() {
		var record UserBotRecord
		if err := rows.Scan(&record.ID, &record.UserID, &record.BotID); err != nil {
			return nil, err
		}
		results = append(results, record.ToModel())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, store.ErrUserBotNotFound
	}
	return results, nil
}
