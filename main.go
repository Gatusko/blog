package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type resp struct {
	Message string `json:"error"`
}

type ready struct {
	Status string `json:"status"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Issue loading env file %s", err)
		return
	}
	port := os.Getenv("PORT")
	log.Printf("Running on port : %s", port)
	route := chi.NewRouter()
	route.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))
	route.Mount("/v1", routeHandler())
	err = http.ListenAndServe(":"+port, route)
	log.Fatalf("Got an error running server: %s", err)
}

func routeHandler() http.Handler {
	route := chi.NewRouter()
	route.Get("/health", healthCheckHandler)
	route.Get("/error", errorResponseHandler)
	return route
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := ready{"ok"}
	respondWithJSON(w, 200, status)
}

func errorResponseHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "error testing")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json : %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("SErver error 5xx: %s", msg)
	}
	res := resp{msg}
	respondWithJSON(w, code, res)
}
