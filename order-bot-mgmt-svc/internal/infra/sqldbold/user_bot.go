package sqldbold

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldbold/pqsqldb"
	"order-bot-mgmt-svc/internal/infra/sqldbold/sqldbexecutor"
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
	insertUserBotQuery   = `INSERT INTO user_bot (id, user_id, bot_id) VALUES ($1, $2, $3);`
	selectUserBotsByUser = `SELECT id, user_id, bot_id FROM user_bot WHERE user_id = $1;`
)

func NewUserBotStore(db *pqsqldb.DB) *UserBotStore {
	if db == nil {
		panic("sqldb.NewUserBotStore(), the db ptr is nil")
	}
	return &UserBotStore{db: db.Conn()}
}

func (s *UserBotStore) Create(ctx context.Context, tx store.Tx, userBot entities.UserBot) error {
	record := UserBotRecordFromModel(userBot)
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.UserBotStore.Create(): %w", err)
	}
	_, err = exec.ExecContext(ctx, insertUserBotQuery, record.ID, record.UserID, record.BotID)
	if err != nil {
		return fmt.Errorf("sqldb.UserBotStore.Create(), ExecContext: %w", err)
	}
	return nil
}

func (s *UserBotStore) FindByUserID(ctx context.Context, tx store.Tx, userID string) ([]entities.UserBot, error) {
	exec, err := sqldbexecutor.Executor(s.db, tx)
	if err != nil {
		return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", err)
	}
	rows, err := exec.QueryContext(ctx, selectUserBotsByUser, userID)
	if err != nil {
		return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", err)
	}
	defer rows.Close()
	var results []entities.UserBot
	for rows.Next() {
		var record UserBotRecord
		if err := rows.Scan(&record.ID, &record.UserID, &record.BotID); err != nil {
			return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", err)
		}
		results = append(results, record.ToModel())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", store.ErrUserBotNotFound)
	}
	return results, nil
}
