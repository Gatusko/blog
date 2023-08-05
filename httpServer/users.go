package main

import (
	"fmt"
	"github.com/Gatusko/blog/internal/data"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type InputUser struct {
	name string
}

func (app *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, "")
}

func (app *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {

	inputUser := InputUser{}
	err := app.readJSON(w, r, &inputUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("%s", err))
		return
	}
	user, _ := data.NewUser(inputUser.name)
	err = app.models.Users.Insert(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%s", err))
		return
	}
	respondWithJSON(w, 200, user)
}

func (app *apiConfig) userRoute() http.Handler {
	route := chi.NewRouter()
	route.Get("/", app.getUsers)
	route.Post("/", app.postUsers)
	return route
}
