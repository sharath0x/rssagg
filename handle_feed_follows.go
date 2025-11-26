package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sharath0x/rssagg/internal/database"
)

func (apiCfg *apiConfig) handleCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	type Parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	params := Parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("cannot parse the json, %v", err))
		return
	}
	fmt.Println(params.FeedID)
	feed_follows, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to create feed_follows, %v", err))
		return
	}

	respondWithJson(w, 200, DatabaseFeedFollowsToFeedFollows(feed_follows))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedfollows, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to fetch the feed follows, %v", err))
		return
	}
	respondWithJson(w, 200, DatabaseFeedsFollowsToFeedsFollows(feedfollows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowsId")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to parse url, %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to delete the feedfollow, %v", err))
		return
	}
	respondWithJson(w, 200, struct{}{})
}
