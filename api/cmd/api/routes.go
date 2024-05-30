package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	app.authRoutes(router)
	app.checkInsRoutes(router)
	app.locationsRoutes(router)
	app.schoolsRoutes(router)
	app.usersRoutes(router)
	app.websocketsRoutes(router)

	var sentryClientOptions *sentry.ClientOptions = nil
	if len(app.config.SentryDsn) > 0 {
		sentryClientOptions = &sentry.ClientOptions{
			Dsn:              app.config.SentryDsn,
			Environment:      app.config.Env,
			Release:          app.config.Release,
			EnableTracing:    true,
			TracesSampleRate: app.config.SampleRate,
		}
	}

	allowedOrigins := []string{app.config.WebURL}
	handlers := middleware.Default(app.config.Throttle, allowedOrigins, sentryClientOptions)

	standard := alice.New(handlers...)

	return standard.Then(router)
}
