package sqldb

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type MenuRecord struct {
	ID    string `gorm:"column:id;primaryKey"`
	BotID string `gorm:"column:bot_id"`
}

func (MenuRecord) TableName() string { return "menu" }

func MenuRecordFromModel(menu entities.Menu) MenuRecord {
	return MenuRecord{ID: menu.ID, BotID: menu.BotID}
}
func (r MenuRecord) ToModel() entities.Menu { return entities.Menu{ID: r.ID, BotID: r.BotID} }

type MenuStore struct{ db *gorm.DB }

func NewMenuStore(db *DB) *MenuStore {
	if db == nil {
		panic("sqldb.NewMenuStore(), the db ptr is nil")
	}
	return &MenuStore{db: db.Gorm()}
}

func (s *MenuStore) FindByBotID(ctx context.Context, botID string) (entities.Menu, error) {
	var record MenuRecord
	if err := s.db.WithContext(ctx).Where("bot_id = ?", botID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Menu{}, fmt.Errorf("sqldb.MenuStore.FindByBotID: %w", store.ErrMenuNotFound)
		}
		return entities.Menu{}, fmt.Errorf("sqldb.MenuStore.FindByBotID: %w", err)
	}
	return record.ToModel(), nil
}
func (s *MenuStore) CreateMenu(ctx context.Context, tx store.Tx, menu entities.Menu) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.CreateMenu: %w", err)
	}
	record := MenuRecordFromModel(menu)
	if err := db.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("sqldb.MenuStore.CreateMenu: %w", err)
	}
	return nil
}
func (s *MenuStore) UpdateMenu(ctx context.Context, tx store.Tx, menu entities.Menu) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu: %w", err)
	}
	res := db.WithContext(ctx).Model(&MenuRecord{}).Where("id = ?", menu.ID).Update("bot_id", menu.BotID)
	if res.Error != nil {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("sqldb.MenuStore.UpdateMenu: %w", store.ErrMenuNotFound)
	}
	return nil
}
func (s *MenuStore) DeleteMenu(ctx context.Context, tx store.Tx, menuID string) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", err)
	}
	res := db.WithContext(ctx).Where("id = ?", menuID).Delete(&MenuRecord{})
	if res.Error != nil {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("sqldb.MenuStore.DeleteMenu: %w", store.ErrMenuNotFound)
	}
	return nil
}
