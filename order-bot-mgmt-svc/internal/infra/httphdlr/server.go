package httphdlr

import (
	"net/http"
	"order-bot-mgmt-svc/internal/services"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/store"
)

type ServerContainer struct {
	port int

	idb      store.DB
	services *services.Services
}

func NewServerContainer(port int, db store.DB, services *services.Services) *ServerContainer {
	return &ServerContainer{port: port, idb: db, services: services}
}

func (s *ServerContainer) db() store.DB              { return s.idb }
func (s *ServerContainer) AuthService() *authsvc.Svc { return s.services.Auth }
func (s *ServerContainer) MenuService() *menusvc.Svc { return s.services.Menu }
func (s *ServerContainer) BotService() *botsvc.Svc   { return s.services.Bot }
func (s *ServerContainer) OrderService() *ordersvc.Svc {
	return s.services.Order
}

func NewHTTPServer(s *ServerContainer, ginMode string, addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: NewRouters(s, ginMode),
	}
}
