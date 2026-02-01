package httphdlrs

import (
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/util"
)

type menuReq struct {
	BotID string        `json:"bot_id" validate:"required"`
	Items []menuItemReq `json:"items"`
}

type menuRes struct {
	BotID string        `json:"bot_id""`
	Items []menuItemRes `json:"items"`
}

type menuItemReq struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type menuItemRes struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func modelFromMenReq(req menuReq) []entities.MenuItem {
	items := make([]entities.MenuItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, entities.MenuItem{
			ID:           util.NewID(),
			MenuItemName: item.Name,
			Price:        item.Price,
		})
	}
	return items
}

func menuResFromModel(menu entities.Menu, items []entities.MenuItem) menuRes {
	resItems := make([]menuItemRes, 0, len(items))
	for _, item := range items {
		resItems = append(resItems, menuItemRes{
			ID:    item.ID,
			Name:  item.MenuItemName,
			Price: item.Price,
		})
	}
	return menuRes{
		BotID: menu.BotID,
		Items: resItems,
	}
}
