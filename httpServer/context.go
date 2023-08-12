package main

import (
	"context"
	"github.com/Gatusko/blog/internal/data"
	"net/http"
)

const userContextKey = "user"

func (api *apiConfig) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (api *apiConfig) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("Missin user value in request context")
	}
	return user
}
