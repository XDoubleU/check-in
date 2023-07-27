package main

import (
	"net/http"

	"github.com/goddtriffin/helmet"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	helmet := helmet.Default()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{app.config.WebURL},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
	})

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
		cors.Handler,
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
