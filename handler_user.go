package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kuangyuwu/blog-aggregator-bootdev/internal/database"
)

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

	payload := struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Name      string `json:"name"`
		ApiKey    string `json:"api_key"`
	}{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}

	respondWithJSON(w, http.StatusCreated, payload)
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

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, _ := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Error getting user")
		return
	}

	payload := struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Name      string `json:"name"`
		ApiKey    string `json:"api_key"`
	}{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}

	respondWithJSON(w, http.StatusOK, payload)
}
