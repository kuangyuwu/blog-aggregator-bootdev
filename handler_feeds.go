package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kuangyuwu/blog-aggregator-bootdev/internal/database"
)

func (cfg *apiConfig) authedHandlerCreateFeed(w http.ResponseWriter, r *http.Request, u database.User) {
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    u.ID,
	})
	if err != nil {
		log.Printf("Error creating feed: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating feed")
		return
	}

	payload := struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Name      string `json:"name"`
		Url       string `json:"url"`
		UserID    string `json:"user_id"`
	}{
		ID:        feed.ID.String(),
		CreatedAt: feed.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: feed.UpdatedAt.UTC().Format(time.RFC3339),
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID.String(),
	}

	respondWithJSON(w, http.StatusOK, payload)
}
