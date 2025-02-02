package main

import (
	"net/http"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/context"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func (app *Application) checkInsRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /checkins/schools",
		app.authAccess(defaultRole, app.getSortedSchoolsHandler),
	)
	mux.HandleFunc(
		"POST /checkins",
		app.authAccess(defaultRole, app.createCheckInHandler),
	)
}

//	@Summary	Get all schools sorted based on checkins at location
//	@Tags		checkins
//	@Success	200	{object}	[]School
//	@Failure	401	{object}	ErrorDto
//	@Failure	500	{object}	ErrorDto
//	@Router		/checkins/schools [get]

func (app *Application) getSortedSchoolsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	schools, err := app.services.CheckInsWriter.GetAllSchoolsSortedByLocation(
		r.Context(),
		user,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, schools, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Create check-in at location of logged in user
// @Tags		checkins
// @Param		createCheckInDto	body		CreateCheckInDto	true	"CreateCheckInDto"
// @Success	201					{object}	CheckInDto
// @Failure	400					{object}	ErrorDto
// @Failure	401					{object}	ErrorDto
// @Failure	404					{object}	ErrorDto
// @Failure	500					{object}	ErrorDto
// @Router		/checkins [post].
func (app *Application) createCheckInHandler(w http.ResponseWriter, r *http.Request) {
	var createCheckInDto dtos.CreateCheckInDto

	err := httptools.ReadJSON(r.Body, &createCheckInDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	if v, validationErrors := createCheckInDto.Validate(); !v {
		httptools.FailedValidationResponse(w, r, validationErrors)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	checkInDto, err := app.services.CheckInsWriter.Create(
		r.Context(),
		createCheckInDto,
		user,
	)
	if err != nil {
		httptools.HandleError(w, r, err, nil)
		return
	}

	err = httptools.WriteJSON(w, http.StatusCreated, checkInDto, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
