package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/services"
	"check-in/api/internal/validator"
)

func (app *application) locationsRoutes(router *httprouter.Router) {
	router.HandlerFunc(
		http.MethodGet,
		"/locations/:id/checkins/range",
		app.authAccess(allRoles, app.getLocationCheckInsRangeHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/locations/:id/checkins/day",
		app.authAccess(allRoles, app.getLocationCheckInsDayHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/locations/:id",
		app.authAccess(allRoles, app.getLocationHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/locations",
		app.authAccess(managerAndAdminRole, app.getPaginatedLocationsHandler),
	)
	router.HandlerFunc(
		http.MethodPost,
		"/locations",
		app.authAccess(managerAndAdminRole, app.createLocationHandler),
	)
	router.HandlerFunc(
		http.MethodPatch,
		"/locations/:id",
		app.authAccess(allRoles, app.updateLocationHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/locations/:id",
		app.authAccess(managerAndAdminRole, app.deleteLocationHandler),
	)
}

//	@Summary	Get all check-ins at location for a specified day in a specified format
//	@Tags		locations
//	@Param		id			path		string	true	"Location ID"
//	@Param		returnType	query		string	true	"ReturnType ('raw' or 'csv')"
//	@Param		date		query		string	true	"Date (format: 'yyyy-mm-dd')"
//	@Success	200			{object}	[]CheckInsLocationEntryRaw
//	@Failure	400			{object}	ErrorDto
//	@Failure	401			{object}	ErrorDto
//	@Failure	404			{object}	ErrorDto
//	@Failure	500			{object}	ErrorDto
//	@Router		/locations/{id}/checkins/day [get].
func (app *application) getLocationCheckInsDayHandler(w http.ResponseWriter,
	r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	returnType := helpers.ReadStrQueryParam(r, "returnType", "")
	if returnType == "" {
		app.badRequestResponse(w, r, errors.New("missing returnType param in query"))
		return
	}

	date, err := helpers.ReadDateQueryParam(r, "date", nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if date == nil {
		app.badRequestResponse(w, r, errors.New("missing date param in query"))
		return
	}

	user := app.contextGetUser(r)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		app.notFoundResponse(w, r, err, "location", "id", id)
		return
	}

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	startDate := helpers.StartOfDay(date)
	endDate := helpers.EndOfDay(date)

	checkIns, err := app.services.CheckIns.GetAllInRange(
		r.Context(),
		location.ID,
		startDate,
		endDate,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	checkInEntries := app.services.Locations.GetCheckInsEntriesDay(
		checkIns,
		schools,
	)

	if returnType == "csv" {
		filename := time.Now().Format("yyyyMMddHHmmss")

		data := dtos.ConvertCheckInsLocationEntryRawMapToCsv(checkInEntries)
		err = helpers.WriteCSV(w, filename, data)
	} else {
		err = helpers.WriteJSON(w, http.StatusOK, checkInEntries, nil)
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//	@Summary	Get all check-ins at location for a specified range in a specified format
//	@Tags		locations
//	@Param		id			path		string	true	"Location ID"
//	@Param		returnType	query		string	true	"ReturnType ('raw' or 'csv')"
//	@Param		startDate	query		string	true	"StartDate (format: 'yyyy-mm-dd')"
//	@Param		endDate		query		string	true	"EndDate (format: 'yyyy-mm-dd')"
//	@Success	200			{object}	[]CheckInsLocationEntryRaw
//	@Failure	400			{object}	ErrorDto
//	@Failure	401			{object}	ErrorDto
//	@Failure	404			{object}	ErrorDto
//	@Failure	500			{object}	ErrorDto
//	@Router		/locations/{id}/checkins/range [get].
func (app *application) getLocationCheckInsRangeHandler(w http.ResponseWriter,
	r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	returnType := helpers.ReadStrQueryParam(r, "returnType", "")
	if returnType == "" {
		app.badRequestResponse(w, r, errors.New("missing returnType param in query"))
		return
	}

	startDate, err := helpers.ReadDateQueryParam(r, "startDate", nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if startDate == nil {
		app.badRequestResponse(w, r, errors.New("missing startDate param in query"))
		return
	}

	endDate, err := helpers.ReadDateQueryParam(r, "endDate", nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if endDate == nil {
		app.badRequestResponse(w, r, errors.New("missing endDate param in query"))
		return
	}

	user := app.contextGetUser(r)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		app.notFoundResponse(w, r, err, "location", "id", id)
		return
	}

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	startDate = helpers.StartOfDay(startDate)
	endDate = helpers.EndOfDay(endDate)

	checkIns, err := app.services.CheckIns.GetAllInRange(
		r.Context(),
		location.ID,
		startDate,
		endDate,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	checkInEntries := app.services.Locations.GetCheckInsEntriesRange(
		startDate,
		endDate,
		checkIns,
		schools,
	)

	if returnType == "csv" {
		filename := time.Now().Format("yyyyMMddHHmmss")

		data := dtos.ConvertCheckInsLocationEntryRawMapToCsv(checkInEntries)
		err = helpers.WriteCSV(w, filename, data)
	} else {
		err = helpers.WriteJSON(w, http.StatusOK, checkInEntries, nil)
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//	@Summary	Get single location
//	@Tags		locations
//	@Param		id	path		string	true	"Location ID"
//	@Success	200	{object}	models.Location
//	@Failure	400	{object}	ErrorDto
//	@Failure	401	{object}	ErrorDto
//	@Failure	404	{object}	ErrorDto
//	@Failure	500	{object}	ErrorDto
//	@Router		/locations/{id} [get].
func (app *application) getLocationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		app.notFoundResponse(w, r, err, "location", "id", id)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//	@Summary	Get all locations paginated
//	@Tags		locations
//	@Param		page	query		int	false	"Page to fetch"
//	@Success	200		{object}	dtos.PaginatedLocationsDto
//	@Failure	400		{object}	ErrorDto
//	@Failure	401		{object}	ErrorDto
//	@Failure	500		{object}	ErrorDto
//	@Router		/locations [get].
func (app *application) getPaginatedLocationsHandler(w http.ResponseWriter,
	r *http.Request) {
	var pageSize int64 = 3

	page, err := helpers.ReadIntQueryParam(r, "page", 1)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	result, err := getAllPaginated[models.Location](
		r.Context(),
		app.services.Locations,
		page,
		pageSize,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//	@Summary	Create location
//	@Tags		locations
//	@Param		createLocationDto	body		CreateLocationDto	true	"CreateLocationDto"
//	@Success	201					{object}	models.Location
//	@Failure	400					{object}	ErrorDto
//	@Failure	401					{object}	ErrorDto
//	@Failure	409					{object}	ErrorDto
//	@Failure	500					{object}	ErrorDto
//	@Router		/locations [post].
func (app *application) createLocationHandler(w http.ResponseWriter, r *http.Request) {
	var createLocationDto dtos.CreateLocationDto

	err := helpers.ReadJSON(r.Body, &createLocationDto)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if dtos.ValidateCreateLocationDto(v, createLocationDto); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	existingLocation, err := app.services.Locations.GetByName(
		r.Context(),
		createLocationDto.Name,
	)
	if existingLocation != nil || !errors.Is(err, services.ErrRecordNotFound) {
		app.conflictResponse(w, r, err, "location", "name", createLocationDto.Name)
		return
	}

	existingUser, err := app.services.Users.GetByUsername(
		r.Context(),
		createLocationDto.Username,
	)
	if existingUser != nil || !errors.Is(err, services.ErrRecordNotFound) {
		app.conflictResponse(w, r, err, "user", "username", createLocationDto.Username)
		return
	}

	user, err := app.services.Users.Create(
		r.Context(),
		createLocationDto.Username,
		createLocationDto.Password,
		models.DefaultRole,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	location, err := app.services.Locations.Create(
		r.Context(),
		createLocationDto.Name,
		createLocationDto.Capacity,
		user.ID,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusCreated, location, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//	@Summary	Update location
//	@Tags		locations
//	@Param		id					path		string				true	"Location ID"
//	@Param		updateLocationDto	body		UpdateLocationDto	true	"UpdateLocationDto"
//	@Success	200					{object}	models.Location
//	@Failure	400					{object}	ErrorDto
//	@Failure	401					{object}	ErrorDto
//	@Failure	409					{object}	ErrorDto
//	@Failure	500					{object}	ErrorDto
//	@Router		/locations/{id} [patch].
func (app *application) updateLocationHandler(w http.ResponseWriter,
	r *http.Request) {
	var updateLocationDto dtos.UpdateLocationDto

	id, err := helpers.ReadUUIDURLParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = helpers.ReadJSON(r.Body, &updateLocationDto)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if dtos.ValidateUpdateLocationDto(v, updateLocationDto); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	hasConflicts := app.checkForConflictsOnUpdate(w, r, updateLocationDto)
	if hasConflicts {
		return
	}

	user := app.contextGetUser(r)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		app.notFoundResponse(w, r, err, "location", "id", id)
		return
	}

	locationUser, err := app.services.Users.GetByID(
		r.Context(),
		location.UserID,
		models.DefaultRole,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.services.Locations.Update(
		r.Context(),
		location,
		locationUser,
		updateLocationDto,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.services.WebSockets.AddUpdateEvent(*location)

	err = helpers.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) checkForConflictsOnUpdate(
	w http.ResponseWriter,
	r *http.Request,
	updateLocationDto dtos.UpdateLocationDto,
) bool {
	if updateLocationDto.Name != nil {
		existingLocation, err := app.services.Locations.GetByName(
			r.Context(),
			*updateLocationDto.Name,
		)

		if existingLocation != nil || !errors.Is(err, services.ErrRecordNotFound) {
			app.conflictResponse(w, r, err, "location", "name", *updateLocationDto.Name)
			return true
		}
	}

	if updateLocationDto.Username != nil {
		existingUser, err := app.services.Users.GetByUsername(
			r.Context(),
			*updateLocationDto.Username,
		)

		if existingUser != nil || !errors.Is(err, services.ErrRecordNotFound) {
			app.conflictResponse(
				w,
				r,
				err,
				"user",
				"username",
				*updateLocationDto.Username,
			)
			return true
		}
	}

	return false
}

//	@Summary	Delete location
//	@Tags		locations
//	@Param		id	path		string	true	"Location ID"
//	@Success	200	{object}	models.Location
//	@Failure	400	{object}	ErrorDto
//	@Failure	401	{object}	ErrorDto
//	@Failure	404	{object}	ErrorDto
//	@Failure	500	{object}	ErrorDto
//	@Router		/locations/{id} [delete].
func (app *application) deleteLocationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil {
		app.notFoundResponse(w, r, err, "location", "id", id)
		return
	}

	err = app.services.Locations.Delete(r.Context(), location.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
