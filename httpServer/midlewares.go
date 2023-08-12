package main

import (
	"github.com/Gatusko/blog/internal/data"
	"net/http"
	"strings"
)

type authedHandler func(http.ResponseWriter, *http.Request, data.User)

func (api *apiConfig) middlewareAuth(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
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
		api.contextSetUser(r, &user)
		next.ServeHTTP(w, r)
	})
}
