package middleware

import (
	"errors"
	"net/http"

	"github.com/Galish/loyalty-system/internal/app/services/auth"
	"github.com/Galish/loyalty-system/internal/logger"
)

var errMissingUserID = errors.New("user id not specified")

func WithAuthChecker(authManager auth.AuthManager) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(auth.AuthCookieName)
			if err != nil {
				logger.WithError(err).Debug("unable to extract auth cookie")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			authData, err := authManager.ParseToken(cookie.Value)
			if err != nil {
				logger.WithError(err).Debug("unable to parse auth token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if authData.UserID == "" {
				logger.WithError(errMissingUserID).Debug("unauthorized access attempt")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r.Header.Set(auth.AuthHeaderName, authData.UserID)

			h.ServeHTTP(w, r)
		})
	}
}
