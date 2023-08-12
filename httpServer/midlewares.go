package main

import (
	"github.com/Gatusko/blog/internal/data"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, data.User)

func (api *apiConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			respondWithError(w, http.StatusBadRequest, "Missing Authorization Header")
			return
		}
		splitHeader := strings.Split(authorizationHeader, " ")
		if len(splitHeader) != 2 || splitHeader[0] != "ApiKey" {
			respondWithError(w, http.StatusUnauthorized, "Invalid ApiKey")
			return
		}
		user, err := api.models.Users.Get(splitHeader[1])
		if err != nil {
			respondWithError(w, http.StatusNotFound, "ApiKey not found")
			return
		}
		r = api.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
