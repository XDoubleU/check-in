package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"check-in/api/internal/config"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	app.authRoutes(router)
	app.checkInsRoutes(router)
	app.locationsRoutes(router)
	app.schoolsRoutes(router)
	app.usersRoutes(router)
	app.websocketsRoutes(router)

	var sentryClientOptions *sentry.ClientOptions
	if len(app.config.SentryDsn) > 0 {
		sentryClientOptions = &sentry.ClientOptions{
			Dsn:              app.config.SentryDsn,
			Environment:      app.config.Env,
			Release:          app.config.Release,
			EnableTracing:    true,
			TracesSampleRate: app.config.SampleRate,
		}
	}

	isTestEnv := app.config.Env == config.TestEnv
	allowedOrigins := []string{app.config.WebURL}
	handlers := middleware.Default(
		isTestEnv,
		allowedOrigins,
		sentryClientOptions,
		app.config.Env == config.DevEnv || app.config.Env == config.TestEnv,
	)

	standard := alice.New(handlers...)

	return standard.Then(router)
}
