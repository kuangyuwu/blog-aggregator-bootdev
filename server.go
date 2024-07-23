package main

import "net/http"

func (cfg *config) initServer() *http.Server {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", cfg.handlerReadiness)
	mux.HandleFunc("GET /v1/err", cfg.handlerError)

	return &http.Server{
		Addr:    ":" + cfg.port,
		Handler: mux,
	}
}
