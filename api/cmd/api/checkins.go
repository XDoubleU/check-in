package main

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/validator"
)

func (app *application) checkInsRoutes(router *httprouter.Router) {
	router.HandlerFunc(
		http.MethodGet,
		"/checkins/schools",
		app.authAccess(defaultRole, app.getSortedSchoolsHandler),
	)
	router.HandlerFunc(
		http.MethodPost,
		"/checkins",
		app.authAccess(defaultRole, app.createCheckInHandler),
	)
}

//	@Summary	Get all schools sorted based on checkins at location
//	@Tags		checkins
//	@Success	200	{object}	[]School
//	@Failure	401	{object}	ErrorDto
//	@Failure	500	{object}	ErrorDto
//	@Router		/checkins/schools [get]

func (app *application) getSortedSchoolsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	user := app.contextGetUser(r)
	location, err := app.services.Locations.GetByUserID(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schools, err := app.services.Schools.GetAllSortedByLocation(
		r.Context(),
		location.ID,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, schools, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary	Create check-in at location of logged in user
// @Tags		checkins
// @Param		checkInDto	body		CheckInDto	true	"CheckInDto"
// @Success	201			{object}	CheckIn
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	404			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/checkins [post].
func (app *application) createCheckInHandler(w http.ResponseWriter, r *http.Request) {
	var checkInDto dtos.CheckInDto

	err := helpers.ReadJSON(r.Body, &checkInDto)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if dtos.ValidateCheckInDto(v, checkInDto); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)
	location, err := app.services.Locations.GetByUserID(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	school, err := app.services.Schools.GetByID(r.Context(), checkInDto.SchoolID)
	if err != nil {
		app.notFoundResponse(w, r, err, "school", "id", checkInDto.SchoolID, "schoolId")
		return
	}

	if location.Available <= 0 {
		app.badRequestResponse(w, r, errors.New("location has no available spots"))
		return
	}

	checkIn, err := app.services.CheckIns.Create(
		r.Context(),
		location,
		school,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.services.WebSockets.AddUpdateEvent(*location)

	err = helpers.WriteJSON(w, http.StatusCreated, checkIn, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
