package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *apiConfig) mapAllRouters() http.Handler {
	route := chi.NewRouter()

	route.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	route.Mount("/v1/users", app.userRoute())
	route.Mount("/v1", app.healthCheckRoute())
	return route
}
