package httphdlrs

import (
	"errors"
	"log/slog"
	"net/http"
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
	mux.HandleFunc("GET /{botId}", getMenuHdlrFunc(s))
	mux.HandleFunc("PUT /", updateMenuHdlrFunc(s))
	return mux
}

func createMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		req, ok := decodeJsonRequest[menuReq](w, r)
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
		writeJSON(w, http.StatusCreated, menuResFromModel(menu, items))
	}
}

func getMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		botId := r.PathValue("botId")
		menu, items, err := service.GetMenu(botId)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, menuResFromModel(menu, items))
	}
}

func updateMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		req, ok := decodeJsonRequest[menuReq](w, r)
		if !ok {
			return
		}
		if err := validatorutil.RequiredStrings(req); err != nil {
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody)
			return
		}
		menu, items, err := service.UpdateMenu(req.BotID, extractItemNames(req.Items))
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, menuResFromModel(menu, items))
	}
}

func extractItemNames(items []menuItemReq) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
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
