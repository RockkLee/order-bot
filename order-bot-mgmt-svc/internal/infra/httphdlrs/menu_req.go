package httphdlrs

import "order-bot-mgmt-svc/internal/models/entities"

type menuItemReq struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price"`
}

type menuReq struct {
	BotID string        `json:"bot_id" validate:"required"`
	Items []menuItemReq `json:"items"`
}

type menuItemRes struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type menuRes struct {
	ID    string        `json:"id"`
	BotID string        `json:"bot_id"`
	Items []menuItemRes `json:"items"`
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
		ID:    menu.ID,
		BotID: menu.BotID,
		Items: resItems,
	}
}
