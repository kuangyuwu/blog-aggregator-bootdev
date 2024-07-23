package main

import (
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	port string
}

func initConfig() *config {
	godotenv.Load()
	port := os.Getenv("PORT")

	return &config{
		port: port,
	}
}
