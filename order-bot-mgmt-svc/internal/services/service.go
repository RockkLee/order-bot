package services

import (
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/services/ordersvc"

	"order-bot-mgmt-svc/internal/services/authsvc"
)

type Services struct {
	Auth  *authsvc.Svc
	Menu  *menusvc.Svc
	Bot   *botsvc.Svc
	Order *ordersvc.Svc
}

func NewServices(
	auth *authsvc.Svc,
	menu *menusvc.Svc,
	bot *botsvc.Svc,
	order *ordersvc.Svc,
) *Services {
	return &Services{
		Auth:  auth,
		Menu:  menu,
		Bot:   bot,
		Order: order,
	}
}
