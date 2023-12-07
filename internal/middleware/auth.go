package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/logger"
)

var errMissingUserID = errors.New("user id not specified")

func WithAuthChecker(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authData, err := parseAuthCookie(r)
		if err != nil {
			logger.WithError(err).Debug("unauthorized access attempt")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set(auth.AuthHeaderName, authData.UserID)

		h.ServeHTTP(w, r)
	})
}

func parseAuthCookie(r *http.Request) (*auth.JWTClaims, error) {
	cookie, err := r.Cookie(auth.AuthCookieName)
	if err != nil {
		return nil, fmt.Errorf("unable to extract auth cookie: %w", err)
	}

	claims, err := auth.ParseToken(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("unable to parse auth token: %w", err)
	}
	if claims.UserID == "" {
		return nil, errMissingUserID
	}

	return claims, nil
}
