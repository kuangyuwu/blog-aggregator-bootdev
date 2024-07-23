package main

import "log"

func main() {
	cfg := initConfig()
	srv := cfg.initServer()

	log.Printf("Serving on port: %s\n", cfg.port)
	log.Fatal(srv.ListenAndServe())
}
