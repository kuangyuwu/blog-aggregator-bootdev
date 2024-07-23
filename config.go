package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kuangyuwu/blog-aggregator-bootdev/internal/database"
)

type apiConfig struct {
	DB   *database.Queries
	port string
}

func initApiConfig() *apiConfig {

	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	dbQueries := database.New(db)

	return &apiConfig{
		DB:   dbQueries,
		port: port,
	}
}
