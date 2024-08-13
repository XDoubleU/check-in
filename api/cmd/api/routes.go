package main

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/pkg/middleware"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	app.authRoutes(mux)
	app.checkInsRoutes(mux)
	app.locationsRoutes(mux)
	app.schoolsRoutes(mux)
	app.usersRoutes(mux)
	app.websocketsRoutes(mux)
	app.stateRoutes(mux)

	var sentryClientOptions sentry.ClientOptions
	if len(app.config.SentryDsn) > 0 {
		//nolint:exhaustruct //other fields are optional
		sentryClientOptions = sentry.ClientOptions{
			Dsn:              app.config.SentryDsn,
			Environment:      app.config.Env,
			Release:          app.config.Release,
			EnableTracing:    true,
			TracesSampleRate: app.config.SampleRate,
		}
	}

	allowedOrigins := []string{app.config.WebURL}
	handlers, err := middleware.DefaultWithSentry(
		app.logger,
		allowedOrigins,
		app.config.Env,
		sentryClientOptions,
	)

	if err != nil {
		panic(err)
	}

	standard := alice.New(handlers...)
	return standard.Then(mux)
}
