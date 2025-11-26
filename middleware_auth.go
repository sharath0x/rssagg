package main

import (
	"fmt"
	"net/http"

	"github.com/sharath0x/rssagg/internal/auth"
	"github.com/sharath0x/rssagg/internal/database"
)

type authhandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authhandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		key, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), key)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Failed to get user, %v", err))
		}

		handler(w, r, user)
	}
}
