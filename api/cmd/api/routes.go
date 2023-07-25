package main

import (
	"net/http"

	"github.com/goddtriffin/helmet"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	helmet := helmet.Default()

	sentryHandler := app.getSentryHandler()

	app.authRoutes(router)
	app.checkInsRoutes(router)
	app.locationsRoutes(router)
	app.schoolsRoutes(router)
	app.usersRoutes(router)
	app.websocketsRoutes(router)

	middleware := []alice.Constructor{
		helmet.Secure,
		app.recoverPanic,
		app.enableCORS,
	}

	if app.config.Throttle {
		middleware = append(middleware, app.rateLimit)
	}

	if sentryHandler != nil {
		middleware = append(middleware, sentryHandler.Handle)
		middleware = append(middleware, app.enrichSentryHub)
	}

	standard := alice.New(middleware...)

	return standard.Then(router)
}
