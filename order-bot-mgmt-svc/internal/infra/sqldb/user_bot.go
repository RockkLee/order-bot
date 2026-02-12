package sqldb

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type UserBotRecord struct {
	ID     string `gorm:"column:id;primaryKey"`
	UserID string `gorm:"column:user_id"`
	BotID  string `gorm:"column:bot_id"`
}

func (UserBotRecord) TableName() string { return "user_bot" }

func UserBotRecordFromModel(userBot entities.UserBot) UserBotRecord {
	return UserBotRecord{ID: userBot.ID, UserID: userBot.UserID, BotID: userBot.BotID}
}
func (r UserBotRecord) ToModel() entities.UserBot {
	return entities.UserBot{ID: r.ID, UserID: r.UserID, BotID: r.BotID}
}

type UserBotStore struct{ db *gorm.DB }

func NewUserBotStore(db *DB) *UserBotStore {
	if db == nil {
		panic("sqldb.NewUserBotStore(), the db ptr is nil")
	}
	return &UserBotStore{db: db.Gorm()}
}

func (s *UserBotStore) Create(ctx context.Context, tx store.Tx, userBot entities.UserBot) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.UserBotStore.Create: %w", err)
	}
	if err := db.WithContext(ctx).Create(&UserBotRecordFromModel(userBot)).Error; err != nil {
		return fmt.Errorf("sqldb.UserBotStore.Create: %w", err)
	}
	return nil
}

func (s *UserBotStore) FindByUserID(ctx context.Context, tx store.Tx, userID string) ([]entities.UserBot, error) {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", err)
	}
	var records []UserBotRecord
	if err := db.WithContext(ctx).Where("user_id = ?", userID).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("sqldb.UserBotStore.FindByUserID: %w", store.ErrUserBotNotFound)
	}
	results := make([]entities.UserBot, 0, len(records))
	for _, record := range records {
		results = append(results, record.ToModel())
	}
	return results, nil
}
