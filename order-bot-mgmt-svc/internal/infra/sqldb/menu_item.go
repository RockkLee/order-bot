package sqldb

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type MenuItemRecord struct {
	ID           string  `gorm:"column:id;primaryKey"`
	MenuID       string  `gorm:"column:menu_id"`
	MenuItemName string  `gorm:"column:menu_item_name"`
	Price        float64 `gorm:"column:price"`
}

func (MenuItemRecord) TableName() string { return "menu_item" }

func MenuItemRecordFromModel(item entities.MenuItem) MenuItemRecord {
	return MenuItemRecord{ID: item.ID, MenuID: item.MenuID, MenuItemName: item.MenuItemName, Price: item.Price}
}
func (r MenuItemRecord) ToModel() entities.MenuItem {
	return entities.MenuItem{ID: r.ID, MenuID: r.MenuID, MenuItemName: r.MenuItemName, Price: r.Price}
}

type MenuItemStore struct{ db *gorm.DB }

func NewMenuItemStore(db *DB) *MenuItemStore {
	if db == nil {
		panic("sqldb.NewMenuItemStore(), the db ptr is nil")
	}
	return &MenuItemStore{db: db.Gorm()}
}

func (s *MenuItemStore) FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error) {
	var records []MenuItemRecord
	if err := s.db.WithContext(ctx).Where("menu_id = ?", menuID).Order("id").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("sqldb.MenuItemStore.FindItems: %w", err)
	}
	items := make([]entities.MenuItem, 0, len(records))
	for _, record := range records {
		items = append(items, record.ToModel())
	}
	return items, nil
}
func (s *MenuItemStore) DeleteMenuItems(ctx context.Context, tx store.Tx, menuID string) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.DeleteMenuItems: %w", err)
	}
	if err := db.WithContext(ctx).Where("menu_id = ?", menuID).Delete(&MenuItemRecord{}).Error; err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.DeleteMenuItems: %w", err)
	}
	return nil
}
func (s *MenuItemStore) CreateMenuItems(ctx context.Context, tx store.Tx, items []entities.MenuItem) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.CreateMenuItems: %w", err)
	}
	records := make([]MenuItemRecord, 0, len(items))
	for _, item := range items {
		records = append(records, MenuItemRecordFromModel(item))
	}
	if err := db.WithContext(ctx).Create(&records).Error; err != nil {
		return fmt.Errorf("sqldb.MenuItemStore.CreateMenuItems: %w", err)
	}
	return nil
}
