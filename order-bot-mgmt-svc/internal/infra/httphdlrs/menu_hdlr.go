package httphdlrs

import (
	"errors"
	"log/slog"
	"net/http"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util/errutil"
	"order-bot-mgmt-svc/internal/util/validatorutil"
)

type MenuServer interface {
	MenuService() *menusvc.Svc
}

const MenuPrefix = "/menus"

func MenuHdlr(s MenuServer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", createMenuHdlrFunc(s))
	mux.HandleFunc("GET /{menuID}", getMenuHdlrFunc(s))
	mux.HandleFunc("PUT /{menuID}", updateMenuHdlrFunc(s))
	mux.HandleFunc("DELETE /{menuID}", deleteMenuHdlrFunc(s))
	return mux
}

func createMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		req, ok := decodeJsonRequest[menuRequest](w, r)
		if !ok {
			return
		}
		if err := validatorutil.RequiredStrings(req); err != nil {
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
			return
		}
		menu, items, err := service.CreateMenu(req.BotID, extractItemNames(req.Items))
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, menuResponseFromModel(menu, items))
	}
}

func getMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		menuID := r.PathValue("menuID")
		menu, items, err := service.GetMenu(menuID)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, menuResponseFromModel(menu, items))
	}
}

func updateMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		menuID := r.PathValue("menuID")
		req, ok := decodeJsonRequest[menuRequest](w, r)
		if !ok {
			return
		}
		if err := validatorutil.RequiredStrings(req); err != nil {
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
			return
		}
		menu, items, err := service.UpdateMenu(menuID, req.BotID, extractItemNames(req.Items))
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, menuResponseFromModel(menu, items))
	}
}

func deleteMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		menuID := r.PathValue("menuID")
		if err := service.DeleteMenu(menuID); err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "menu deleted"})
	}
}

func extractItemNames(items []menuItemRequest) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
}

func menuResponseFromModel(menu entities.Menu, items []entities.MenuItem) menuResponse {
	respItems := make([]menuItemResponse, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, menuItemResponse{
			ID:   item.ID,
			Name: item.MenuItemName,
		})
	}
	return menuResponse{
		ID:    menu.ID,
		BotID: menu.BotID,
		Items: respItems,
	}
}

func writeMenuError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, menusvc.ErrInvalidMenu):
		WriteError(w, http.StatusBadRequest, menusvc.ErrInvalidMenu.Error())
	case errors.Is(err, store.ErrMenuNotFound):
		WriteError(w, http.StatusNotFound, store.ErrMenuNotFound.Error())
	default:
		WriteError(w, http.StatusInternalServerError, "menu request failed")
	}
}
