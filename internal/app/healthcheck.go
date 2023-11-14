package app

import (
	"dryve/internal/datastruct"
	"dryve/internal/dto"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// HealthcheckEchoAuth is an authenticated simple echo handler that returns the param and the user's email.
func (app *App) HealthcheckEchoAuth(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(*datastruct.User)
	param, _ := strconv.Atoi(chi.URLParam(r, "param"))
	w.Write([]byte(fmt.Sprintf("param: %d, id: %d, email: %s", param, user.ID, user.Email)))
}

// Healthcheck is a public simple handler for healthcheck.
func (app *App) Healthcheck(w http.ResponseWriter, r *http.Request) {
	EncodeJSONAndSend(w, dto.HealthcheckResponse{
		Status: "OK",
		Errors: []string{},
	})
}
