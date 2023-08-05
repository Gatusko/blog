package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *apiConfig) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := ready{"ok"}
	respondWithJSON(w, 200, status)
}

func (app *apiConfig) errorResponseHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "error testing")
}

func (app *apiConfig) healthCheckRoute() http.Handler {
	route := chi.NewRouter()
	route.Get("/health", app.healthCheckHandler)
	route.Get("/error", app.errorResponseHandler)
	return route
}
