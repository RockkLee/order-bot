package httphdlrsold

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
	mux.HandleFunc("POST /{botId}/publish", publishMenuHdlrFunc(s))
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
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody.Error())
			return
		}
		menu, items, err := service.CreateMenu(r.Context(), req.BotID, modelFromMenReq(req))
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
		menu, items, err := service.GetMenu(r.Context(), botId)
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
			WriteError(w, http.StatusBadRequest, ErrMsgInvalidRequestBody.Error())
			return
		}
		menu, items, err := service.UpdateMenu(r.Context(), req.BotID, modelFromMenReq(req))
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, menuResFromModel(menu, items))
	}
}

func publishMenuHdlrFunc(s MenuServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service := s.MenuService()
		botId := r.PathValue("botId")
		menu, items, err := service.PublishMenu(r.Context(), botId)
		if err != nil {
			slog.Error(errutil.FormatErrChain(err))
			writeMenuError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, menuResFromModel(menu, items))
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
