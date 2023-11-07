package billing

import (
	"net/http"
	"os"

	"github.com/devpies/saas-core/internal/billing/config"
	"github.com/devpies/saas-core/internal/billing/handler"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/devpies/saas-core/pkg/web/mid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// Routes composes routes, middleware and handlers.
func Routes(
	log *zap.Logger,
	shutdown chan os.Signal,
	subscriptionHandler *handler.SubscriptionHandler,
	config config.Config,
) http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://devpie.local:3000", "https://devpie.io"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "BasePath"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	middleware := []web.Middleware{
		mid.Logger(log),
		mid.Errors(log),
		mid.Auth(log, config.Cognito.Region, config.Cognito.SharedUserPoolID),
		mid.Panics(log),
	}

	app := web.NewApp(mux, shutdown, log, middleware...)

	app.Handle(http.MethodPost, "/billing", subscriptionHandler.Create)
	app.Handle(http.MethodGet, "/billing/{tenantID}", subscriptionHandler.GetOne)
	app.Handle(http.MethodPost, "/billing/payment-intent", subscriptionHandler.GetPaymentIntent)
	//app.Handle(http.MethodGet, "/billing/cancel", subscriptionHandler.Cancel)
	//app.Handle(http.MethodGet, "/billing/refund", subscriptionHandler.Refund)

	return app
}
