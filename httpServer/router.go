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
	route.Mount("/v1/feeds", app.feedRoute())
	route.Mount("/v1/feed_follows", app.feedFollowRoutes())
	return route
}
