package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kuangyuwu/blog-aggregator-bootdev/internal/database"
)

type FeedFollow struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserID    string `json:"user_id"`
	FeedID    string `json:"feed_id"`
}

func dbFeedFollowtoFeedFollow(ff database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        ff.ID.String(),
		CreatedAt: ff.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: ff.UpdatedAt.UTC().Format(time.RFC3339),
		UserID:    ff.UserID.String(),
		FeedID:    ff.FeedID.String(),
	}
}

func (cfg *apiConfig) authedHandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, u database.User) {

	decoder := json.NewDecoder(r.Body)
	params := struct {
		FeedId string `json:"feed_id"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	feedId, err := uuid.Parse(params.FeedId)
	if err != nil {
		log.Printf("Invalid feed ID: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid feed ID")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    feedId,
	})
	if err != nil {
		log.Printf("Error creating feed follow: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, dbFeedFollowtoFeedFollow(feedFollow))
}

func (cfg *apiConfig) authedHandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollows, err := cfg.DB.GetAllFeedFollowsByUserId(r.Context(), u.ID)
	if err != nil {
		log.Printf("Error getting feed follows: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting feed follows")
		return
	}

	l := len(feedFollows)
	payload := make([]FeedFollow, l)
	for i, ff := range feedFollows {
		payload[i] = dbFeedFollowtoFeedFollow(ff)
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) authedHandlerDeleteFeedFollowById(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollowIdString := r.PathValue("feedFollowId")

	feedFollowId, err := uuid.Parse(feedFollowIdString)
	if err != nil {
		log.Printf("Invalid feed follow ID: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), feedFollowId)
	if err != nil {
		log.Printf("Error deleting feed follows: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error deleting feed follows")
		return
	}

	respondWithNoContent(w)
}
