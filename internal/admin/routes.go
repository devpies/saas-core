package admin

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/admin/config"
	"github.com/devpies/saas-core/internal/admin/handler"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Routes composes routes, middleware and handlers.
func Routes(
	log *zap.Logger,
	shutdown chan os.Signal,
	assets fs.FS,
	authHandler *handler.AuthHandler,
	webPageHandler *handler.WebPageHandler,
	registrationHandler *handler.RegistrationHandler,
	config config.Config,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(loadSession)

	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(assets))))

	app := web.NewApp(mux, shutdown, log, []web.Middleware{mid.Logger(log), mid.Errors(log), mid.Panics(log)}...)

	// Unauthenticated webpages.
	app.Handle(http.MethodGet, "/", withNoSession()(authHandler.LoginPage))
	app.Handle(http.MethodGet, "/force-new-password", withPasswordChallengeSession()(authHandler.ForceNewPasswordPage))
	app.Handle(http.MethodPost, "/secure-new-password", withNoSession()(authHandler.SetupNewUserWithSecurePassword))
	app.Handle(http.MethodPost, "/authenticate", withNoSession()(authHandler.AuthenticateCredentials))

	// Authenticated webpages.
	app.Handle(http.MethodGet, "/admin", withSession()(webPageHandler.DashboardPage))
	app.Handle(http.MethodGet, "/admin/tenants", withSession()(webPageHandler.TenantsPage))
	app.Handle(http.MethodGet, "/admin/create-tenant", withSession()(webPageHandler.CreateTenantPage))
	app.Handle(http.MethodGet, "/admin/logout", withSession()(authHandler.Logout))
	app.Handle(http.MethodGet, "/*", withSession()(webPageHandler.E404Page))

	app.Handle(http.MethodPost, "/api/send-registration",
		mid.Auth(log, config.Cognito.Region, config.Cognito.UserPoolClientID)(registrationHandler.ProcessRegistration))

	return mux
}
