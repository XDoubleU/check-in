package main

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"

	"check-in/api/internal/models"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(
	r *http.Request,
	user *models.User,
) *http.Request {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		hub.Scope().SetUser(sentry.User{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}

	return user
}
