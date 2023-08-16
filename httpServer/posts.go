package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (api *apiConfig) getPosts(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "5"
	}
	user := api.contextGetUser(r)
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Limit should be an Integer")
		return
	}
	posts, err := api.models.Post.GetPostByUser(user.Id, limitInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	respondWithJSON(w, http.StatusOK, posts)
}

func (api *apiConfig) postRoute() http.Handler {
	r := chi.NewRouter()
	r.With(api.middlewareAuth).Get("/", api.getPosts)
	return r
}
