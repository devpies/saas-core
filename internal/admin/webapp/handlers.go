package webapp

import (
	"net/http"

	"go.uber.org/zap"
)

func (app *WebApp) Login(w http.ResponseWriter, r *http.Request) {
	if err := app.render.Template(w, r, "login", nil); err != nil {
		app.logger.Error("", zap.Error(err))
	}
}

func (app *WebApp) Logout(w http.ResponseWriter, r *http.Request) {}

func (app *WebApp) Dashboard(w http.ResponseWriter, r *http.Request) {}