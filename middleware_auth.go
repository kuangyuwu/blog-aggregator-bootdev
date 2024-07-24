package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/kuangyuwu/blog-aggregator-bootdev/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, _ := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			log.Printf("Authentication failed: %s", err)
			respondWithError(w, http.StatusUnauthorized, "Authentication failed")
			return
		}
		handler(w, r, user)
	}
}
