package main

import (
	"context"

	"github.com/getsentry/sentry-go"
	contexttools "github.com/xdoubleu/essentia/pkg/context"

	"check-in/api/internal/models"
)

const userContextKey = contexttools.Key("user")

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
