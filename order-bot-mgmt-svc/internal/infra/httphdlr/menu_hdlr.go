package httphdlr

import (
	"errors"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/errutil"

	"github.com/gin-gonic/gin"
)

type MenuServer interface {
	MenuService() *menusvc.Svc
}

const MenuPrefix = "/menus"

func RegisterMenuRoutes(r gin.IRoutes, s MenuServer) {
	r.POST("/", createMenuHdlrFunc(s))
	r.GET("/:botId", getMenuHdlrFunc(s))
	r.PUT("/", updateMenuHdlrFunc(s))
	r.POST("/:botId/publish", publishMenuHdlrFunc(s))
}

func createMenuHdlrFunc(s MenuServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req menuReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrMsgInvalidRequestBody})
			return
		}
		menu, items, err := s.MenuService().CreateMenu(c.Request.Context(), req.BotID, modelFromMenReq(req))
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(c, err)
			return
		}
		c.JSON(http.StatusCreated, menuResFromModel(menu, items))
	}
}

func getMenuHdlrFunc(s MenuServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		botID := c.Param("botId")
		menu, items, err := s.MenuService().GetMenu(c.Request.Context(), botID)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(c, err)
			return
		}
		c.JSON(http.StatusOK, menuResFromModel(menu, items))
	}
}

func updateMenuHdlrFunc(s MenuServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req menuReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrMsgInvalidRequestBody})
			return
		}
		menu, items, err := s.MenuService().UpdateMenu(c.Request.Context(), req.BotID, modelFromMenReq(req))
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(c, err)
			return
		}
		c.JSON(http.StatusOK, menuResFromModel(menu, items))
	}
}

func publishMenuHdlrFunc(s MenuServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		botID := c.Param("botId")
		menu, items, err := s.MenuService().PublishMenu(c.Request.Context(), botID)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(c, err)
			return
		}
		c.JSON(http.StatusOK, menuResFromModel(menu, items))
	}
}

func writeMenuError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, menusvc.ErrInvalidMenu):
		c.JSON(http.StatusBadRequest, gin.H{"error": menusvc.ErrInvalidMenu.Error()})
	case errors.Is(err, store.ErrMenuNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": store.ErrMenuNotFound.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "menu request failed"})
	}
}
