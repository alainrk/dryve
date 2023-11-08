package app

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *App) GetFile(w http.ResponseWriter, r *http.Request) {
	rid, _ := strconv.Atoi(chi.URLParam(r, "id"))
	id := uint(rid)

	w.WriteHeader(http.StatusCreated)
	EncodeJSONAndSend(w, map[string]any{"id": id})
}
