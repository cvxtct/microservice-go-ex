package main

import (
	"broker/internal/config"
	"broker/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat(("/ping")))

	// this is just the first try out
	mux.Post("/", handlers.Repo.Broker)

	// endpoint for the gRPC log request
	mux.Post("/log-grpc", handlers.Repo.LogViaGRPC)

	// handle every submission from the frontend
	mux.Post("/handle", handlers.Repo.HandleSubmission)

	return mux
}
