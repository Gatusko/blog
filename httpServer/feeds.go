package main

import (
	"fmt"
	"github.com/Gatusko/blog/internal/data"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type FeedInput struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (api *apiConfig) postFeed(w http.ResponseWriter, r *http.Request) {
	feedInput := FeedInput{}
	err := api.readJSON(w, r, &feedInput)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprint(err))
		return
	}
	user := api.contextGetUser(r)
	feed, _ := data.NewFeed(feedInput.Name, feedInput.URL, user.Id)
	feedAndFollow, err := api.models.Feeds.Insert(feed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, feedAndFollow)
}

func (api *apiConfig) getAllFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := api.models.Feeds.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
}

func (api *apiConfig) feedRoute() http.Handler {
	route := chi.NewRouter()
	route.Get("/", api.getAllFeed)
	route.With(api.middlewareAuth).Post("/", api.postFeed)
	return route
}
