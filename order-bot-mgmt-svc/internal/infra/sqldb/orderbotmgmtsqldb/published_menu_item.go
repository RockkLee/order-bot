package orderbotmgmtsqldb

import (
	"database/sql"
	"math"
	"order-bot-mgmt-svc/internal/models/entities"
	"time"
)

type PublishedMenuItemRecord struct {
	ID          string
	BotID       string
	SKU         string
	Name        string
	Description sql.NullString
	PriceCents  int
	IsAvailable bool
	CreatedAt   time.Time
}

func PublishedMenuItemRecordFromModel(botID string, item entities.MenuItem) PublishedMenuItemRecord {
	priceCents := int(math.Round(item.Price * 100))
	return PublishedMenuItemRecord{
		ID:          item.ID,
		BotID:       botID,
		SKU:         item.ID,
		Name:        item.MenuItemName,
		Description: sql.NullString{},
		PriceCents:  priceCents,
		IsAvailable: true,
		CreatedAt:   time.Now().UTC(),
	}
}
