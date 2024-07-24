package main

import "net/http"

func (cfg *apiConfig) initServer() *http.Server {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", cfg.handlerReadiness)
	mux.HandleFunc("GET /v1/err", cfg.handlerError)

	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.authedHandlerGetUser))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.authedHandlerCreateFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.handlerGetAllFeeds)

	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.authedHandlerCreateFeedFollow))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.authedHandlerGetFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowId}", cfg.middlewareAuth(cfg.authedHandlerDeleteFeedFollowById))

	return &http.Server{
		Addr:    ":" + cfg.port,
		Handler: mux,
	}
}
