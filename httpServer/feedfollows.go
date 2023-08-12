package main

import (
	"fmt"
	"github.com/Gatusko/blog/internal/data"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type FeedFollowInput struct {
	Feed uuid.UUID `json:"feed"`
}

func (api *apiConfig) postFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowInput := FeedFollowInput{}
	err := api.readJSON(w, r, &feedFollowInput)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s", err))
		return
	}
	// Bussines logic for creating feedFollow
	user := api.contextGetUser(r)
	feedFollow, _ := data.NewFeedFollows(user.Id, feedFollowInput.Feed)
	err = api.models.FeedFollows.CreateFeedFollow(*feedFollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, feedFollow)
}

func (api *apiConfig) getFeedFollow(w http.ResponseWriter, r *http.Request) {
	user := api.contextGetUser(r)
	feedFollows, err := api.models.FeedFollows.GetAllFeedFollows(user.Id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, feedFollows)
}

func (api *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowId := chi.URLParam(r, "feedFollowId")
	if feedFollowId == "" {
		respondWithError(w, http.StatusBadRequest, "Need to provide Id")
		return
	}
	id, err := uuid.Parse(feedFollowId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s", err))
	}
	err = api.models.FeedFollows.DeleteFeedFollow(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	respondWithJSON(w, http.StatusAccepted, struct{}{})
}

func (api *apiConfig) feedFollowRoutes() http.Handler {
	r := chi.NewRouter()
	r.With(api.middlewareAuth).Post("/", api.postFeedFollow)
	r.With(api.middlewareAuth).Get("/", api.getFeedFollow)
	r.Delete("/{feedFollowId}", api.deleteFeedFollow)
	return r
}
