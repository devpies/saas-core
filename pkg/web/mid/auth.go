package mid

import (
	"fmt"
	"github.com/devpies/core/pkg/web"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func Auth(log *zap.Logger, region string, userPoolClientID string) web.Middleware {
	// this is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		// create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			err := verifyToken(w, r, region, userPoolClientID)
			if err != nil {
				log.Error(err.Error())
				return web.NewRequestError(err, http.StatusUnauthorized)
			}

			return after(w, r)
		}

		return h
	}

	return f
}

func verifyToken(w http.ResponseWriter, r *http.Request, region string, userPoolClientID string) error {
	authHeader := r.Header.Get("Authorization")
	splitAuthHeader := strings.Split(authHeader, " ")

	if len(splitAuthHeader) != 2 {
		return fmt.Errorf("missing or invalid authorization header")
	}

	pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
	formattedURL := fmt.Sprintf(pubKeyURL, region, userPoolClientID)

	keySet, err := jwk.Fetch(r.Context(), formattedURL)
	if err != nil {
		return err
	}

	_, err = jwt.Parse(
		[]byte(splitAuthHeader[1]),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)

	// Add user object to context
	return err
}
