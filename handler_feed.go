package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sharath0x/rssagg/internal/database"
)

func (apiConfig *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type Parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	params := Parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("cant parse the JSON, %v", err))
		return
	}

	feed, err := apiConfig.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to create feed, %v", err))
		return
	}
	respondWithJson(w, 201, DatabaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to get feeds, %v", err))
		return
	}
	respondWithJson(w, 201, DatabaseFeedsToFeeds(feeds))
}
