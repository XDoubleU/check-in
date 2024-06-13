package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/context_tools"
	"github.com/getsentry/sentry-go"

	"check-in/api/internal/models"
)

const userContextKey = context_tools.ContextKey("user")

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

	return context_tools.SetContextValue(r, userContextKey, user)
}
