package main

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/justinas/alice"
	"github.com/xdoubleu/essentia/pkg/middleware"
)

func (app *application) routes() (*http.Handler, error) {
	mux := http.NewServeMux()

	app.authRoutes(mux)
	app.checkInsRoutes(mux)
	app.locationsRoutes(mux)
	app.schoolsRoutes(mux)
	app.usersRoutes(mux)
	app.websocketsRoutes(mux)

	var sentryClientOptions sentry.ClientOptions
	if len(app.config.SentryDsn) > 0 {
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
		return nil, err
	}

	standard := alice.New(handlers...)
	handler := standard.Then(mux)

	return &handler, nil
}
