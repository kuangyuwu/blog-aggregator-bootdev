package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kuangyuwu/blog-aggregator-bootdev/internal/database"
)

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	ApiKey    string `json:"api_key"`
}

func dbUserToUser(u database.User) User {
	return User{
		ID:        u.ID.String(),
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.UTC().Format(time.RFC3339),
		Name:      u.Name,
		ApiKey:    u.ApiKey,
	}
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := struct {
		Name string `json:"name"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	apiKey, err := generateApiKey()
	if err != nil {
		log.Printf("Error generating API key: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating API key")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		ApiKey:    apiKey,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, dbUserToUser(user))
}

func generateApiKey() (string, error) {
	length := 32
	apiKeyBytes := make([]byte, length)
	_, err := rand.Read(apiKeyBytes)
	if err != nil {
		return "", err
	}
	apiKey := hex.EncodeToString(apiKeyBytes)
	return apiKey, nil
}

func (cfg *apiConfig) authedHandlerGetUser(w http.ResponseWriter, r *http.Request, u database.User) {
	respondWithJSON(w, http.StatusOK, dbUserToUser(u))
}
