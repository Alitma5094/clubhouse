package main

import (
	"clubhouse/internal/auth"
	"clubhouse/internal/database"
	"net/http"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
			return
		}
		subject, err := auth.ValidateJWT(token, cfg.JWTSecret, auth.TokenTypeAccess)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
			return
		}

		user, err := cfg.DB.GetUser(r.Context(), subject)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get user")
			return
		}

		handler(w, r, user)

	})
}
