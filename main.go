package main

import (
	"log"

	_ "github.com/lib/pq"
)

func main() {
	cfg := initApiConfig()
	srv := cfg.initServer()

	log.Printf("Serving on port: %s\n", cfg.port)
	log.Fatal(srv.ListenAndServe())
}
