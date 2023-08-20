package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/services"
)

func (app *application) logError(err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter,
	_ *http.Request, status int, message any) {
	env := dtos.ErrorDto{
		Status:  status,
		Error:   http.StatusText(status),
		Message: message,
	}
	err := helpers.WriteJSON(w, status, env, nil)
	if err != nil {
		app.logError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter,
	r *http.Request, err error) {
	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}

	message := "the server encountered a problem and could not process your request"
	if app.config.Env != config.ProdEnv {
		message = err.Error()
	}

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter,
	r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter,
	r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *application) unauthorizedResponse(w http.ResponseWriter,
	r *http.Request, message string) {
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) forbiddenResponse(w http.ResponseWriter,
	r *http.Request) {
	message := "user has no access to this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

func (app *application) conflictResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifier string,
	identifierValue string,
	jsonField string,
) {
	value := helpers.AnyToString(identifierValue)

	if err == nil || errors.Is(err, services.ErrRecordUniqueValue) {
		message := fmt.Sprintf(
			"%s with %s '%s' already exists",
			resourceName,
			identifier,
			value,
		)
		err := make(map[string]string)
		err[jsonField] = message
		app.errorResponse(w, r, http.StatusConflict, err)
	} else {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) notFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	resourceName string,
	identifier string, //nolint:unparam //should keep param
	identifierValue any,
	jsonField string,
) {
	value := helpers.AnyToString(identifierValue)

	if err == nil || errors.Is(err, services.ErrRecordNotFound) {
		message := fmt.Sprintf(
			"%s with %s '%s' doesn't exist",
			resourceName,
			identifier,
			value,
		)

		err := make(map[string]string)
		err[jsonField] = message

		app.errorResponse(w, r, http.StatusNotFound, err)
	} else {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) failedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
