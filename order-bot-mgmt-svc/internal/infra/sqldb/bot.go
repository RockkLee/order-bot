package sqldb

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type BotRecord struct {
	ID      string `gorm:"column:id;primaryKey"`
	BotName string `gorm:"column:bot_name"`
}

func (BotRecord) TableName() string { return "bot" }

func BotRecordFromModel(bot entities.Bot) BotRecord {
	return BotRecord{ID: bot.ID, BotName: bot.BotName}
}
func (r BotRecord) ToModel() entities.Bot { return entities.Bot{ID: r.ID, BotName: r.BotName} }

type BotStore struct{ db *gorm.DB }

func NewBotStore(db *DB) *BotStore {
	if db == nil {
		panic("sqldb.NewBotStore(), the db ptr is nil")
	}
	return &BotStore{db: db.Gorm()}
}

func (s *BotStore) Create(ctx context.Context, tx store.Tx, bot entities.Bot) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.BotStore.Create: %w", err)
	}
	record := BotRecordFromModel(bot)
	if err := db.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("sqldb.BotStore.Create: %w", err)
	}
	return nil
}

func (s *BotStore) FindByID(ctx context.Context, tx store.Tx, id string) (entities.Bot, error) {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return entities.Bot{}, fmt.Errorf("sqldb.BotStore.FindByID: %w", err)
	}
	var record BotRecord
	if err := db.WithContext(ctx).Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Bot{}, fmt.Errorf("sqldb.BotStore.FindByID: %w", store.ErrBotNotFound)
		}
		return entities.Bot{}, fmt.Errorf("sqldb.BotStore.FindByID: %w", err)
	}
	return record.ToModel(), nil
}
