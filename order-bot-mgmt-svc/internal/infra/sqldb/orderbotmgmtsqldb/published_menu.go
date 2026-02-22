package orderbotmgmtsqldb

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type PublishedMenuRecord struct {
	ID    string `gorm:"column:id;primaryKey"`
	BotID string `gorm:"column:bot_id"`
}

func (PublishedMenuRecord) TableName() string { return "published_menu" }

type PublishedMenuItemRecord struct {
	ID           string  `gorm:"column:id;primaryKey"`
	MenuID       string  `gorm:"column:menu_id"`
	MenuItemName string  `gorm:"column:menu_item_name"`
	Price        float64 `gorm:"column:price"`
}

func (PublishedMenuItemRecord) TableName() string { return "published_menu_item" }

type PublishedMenuStore struct{ db *gorm.DB }

func NewPublishedMenuStore(db *sqldb.DB) *PublishedMenuStore {
	if db == nil {
		panic("orderbotmgmtsqldb.NewPublishedMenuStore(), the db ptr is nil")
	}
	return &PublishedMenuStore{db: db.Gorm()}
}

func (s *PublishedMenuStore) IsMenuPublished(ctx context.Context, menuID string) (bool, error) {
	var record PublishedMenuRecord
	err := s.db.WithContext(ctx).Where("id = ?", menuID).Take(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.IsMenuPublished: %w", err)
	}
	return true, nil
}

func (s *PublishedMenuStore) ReplaceMenuItems(ctx context.Context, tx store.Tx, menu entities.Menu, items []entities.MenuItem) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems: %w", err)
	}
	if err := db.WithContext(ctx).Where("menu_id = ?", menu.ID).Delete(&PublishedMenuItemRecord{}).Error; err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), delete menu_item: %w", err)
	}
	if err := db.WithContext(ctx).Where("bot_id = ?", menu.BotID).Delete(&PublishedMenuRecord{}).Error; err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), delete menu: %w", err)
	}
	if err := db.WithContext(ctx).Create(&PublishedMenuRecord{ID: menu.ID, BotID: menu.BotID}).Error; err != nil {
		return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), insert menu: %w", err)
	}
	records := make([]PublishedMenuItemRecord, 0, len(items))
	for _, item := range items {
		records = append(records, PublishedMenuItemRecord{ID: item.ID, MenuID: item.MenuID, MenuItemName: item.MenuItemName, Price: item.Price})
	}
	if len(records) > 0 {
		if err := db.WithContext(ctx).Create(&records).Error; err != nil {
			return fmt.Errorf("orderbotmgmtsqldb.PublishedMenuStore.ReplaceMenuItems(), insert: %w", err)
		}
	}
	return nil
}
