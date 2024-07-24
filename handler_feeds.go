package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kuangyuwu/blog-aggregator-bootdev/internal/database"
)

type Feed struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Url       string `json:"url"`
	UserID    string `json:"user_id"`
}

func dbFeedtoFeed(f database.Feed) Feed {
	return Feed{
		ID:        f.ID.String(),
		CreatedAt: f.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: f.UpdatedAt.UTC().Format(time.RFC3339),
		Name:      f.Name,
		Url:       f.Url,
		UserID:    f.UserID.String(),
	}
}

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

	respondWithJSON(w, http.StatusOK, dbFeedtoFeed(feed))
}

func (cfg *apiConfig) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeed(r.Context())
	if err != nil {
		log.Printf("Error getting feeds: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error getting feeds")
		return
	}

	l := len(feeds)
	payload := make([]Feed, l)
	for i, f := range feeds {
		payload[i] = dbFeedtoFeed(f)
	}

	respondWithJSON(w, http.StatusOK, payload)
}
