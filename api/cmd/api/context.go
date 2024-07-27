package main

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/xdoubleu/essentia/pkg/contexttools"

	"check-in/api/internal/models"
)

const userContextKey = contexttools.ContextKey("user")

func (app *Application) contextSetUser(
	ctx context.Context,
	user models.User,
) context.Context {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		//nolint:exhaustruct //other fields are optional
		hub.Scope().SetUser(sentry.User{
			ID:       user.ID,
			Username: user.Username,
		})
	}

	return context.WithValue(ctx, userContextKey, user)
}
