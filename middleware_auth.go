package main

import (
	"net/http"

	"github.com/SudilHasitha/rss_aggregator/internal/database"
	"github.com/SudilHasitha/rss_aggregator/internal/database/auth"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) authMiddleware(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get user")
			return
		}

		handler(w, r, user)
	}
}
