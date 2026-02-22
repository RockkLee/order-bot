package httphdlr

import (
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/services/ordersvc"
	"order-bot-mgmt-svc/internal/util/errutil"

	"github.com/gin-gonic/gin"
)

type OrderServer interface {
	OrderService() *ordersvc.Svc
}

const OrderPrefix = "/orders"

type orderItemRes struct {
	ID               string `json:"id"`
	OrderID          string `json:"order_id"`
	MenuItemID       string `json:"menu_item_id"`
	Name             string `json:"name"`
	Quantity         int    `json:"quantity"`
	UnitPriceScaled  int    `json:"unit_price_scaled"`
	TotalPriceScaled int    `json:"total_price_scaled"`
}

type orderRes struct {
	ID          string         `json:"id"`
	BotID       string         `json:"bot_id"`
	CartID      string         `json:"cart_id"`
	SessionID   string         `json:"session_id"`
	TotalScaled int            `json:"total_scaled"`
	Items       []orderItemRes `json:"items"`
}

func RegisterOrderRoutes(r gin.IRoutes, s OrderServer) {
	r.GET("/", getOrdersHdlrFunc(s))
}

func getOrdersHdlrFunc(s OrderServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ordersWithItems, err := s.OrderService().GetOrdersWithItems(c.Request.Context())
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load orders"})
			return
		}
		response := make([]orderRes, 0, len(ordersWithItems))
		for _, orderWithItems := range ordersWithItems {
			items := make([]orderItemRes, 0, len(orderWithItems.Items))
			for _, item := range orderWithItems.Items {
				items = append(items, orderItemRes{
					ID:               item.ID,
					OrderID:          item.OrderID,
					MenuItemID:       item.MenuItemID,
					Name:             item.Name,
					Quantity:         item.Quantity,
					UnitPriceScaled:  item.UnitPriceScaled,
					TotalPriceScaled: item.TotalPriceScaled,
				})
			}
			response = append(response, orderRes{
				ID:          orderWithItems.Order.ID,
				BotID:       orderWithItems.Order.BotID,
				CartID:      orderWithItems.Order.CartID,
				SessionID:   orderWithItems.Order.SessionID,
				TotalScaled: orderWithItems.Order.TotalScaled,
				Items:       items,
			})
		}
		c.JSON(http.StatusOK, gin.H{"orders": response})
	}
}
