package main

import (
	"net/http"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"

	"check-in/api/internal/dtos"
)

func (app *Application) stateRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /state",
		app.getStateHandler,
	)
	mux.HandleFunc(
		"PATCH /state",
		app.authAccess(adminRole, app.updateStateHandler),
	)
}

// @Summary	Get current state
// @Tags		state
// @Success	200	{object}	State
// @Failure	500	{object}	ErrorDto
// @Router		/state [get].
func (app *Application) getStateHandler(w http.ResponseWriter, r *http.Request) {
	err := httptools.WriteJSON(w, http.StatusOK, app.services.State.Current.Get(), nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Update state
// @Tags		state
// @Param		stateDto	body		StateDto	true	"StateDto"
// @Success	200			{object}	State
// @Failure	500			{object}	ErrorDto
// @Router		/state [patch].
func (app *Application) updateStateHandler(w http.ResponseWriter, r *http.Request) {
	var stateDto *dtos.StateDto

	err := httptools.ReadJSON(r.Body, &stateDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	state, err := app.services.State.UpdateState(r.Context(), stateDto)
	if err != nil {
		httptools.HandleError(w, r, err, nil)
	}

	err = httptools.WriteJSON(w, http.StatusOK, state, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
