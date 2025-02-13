package main

import (
	"fmt"
	"net/http"
	"strconv"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func (app *Application) locationsRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /all-locations/checkins/range",
		app.authAccess(allRoles, app.getLocationCheckInsRangeHandler),
	)
	mux.HandleFunc(
		"GET /all-locations/checkins/day",
		app.authAccess(allRoles, app.getLocationCheckInsDayHandler),
	)
	mux.HandleFunc(
		"GET /locations/{locationId}/checkins",
		app.authAccess(allRoles, app.getAllCheckInsTodayHandler),
	)
	mux.HandleFunc(
		"DELETE /locations/{locationId}/checkins/{checkInId}",
		app.authAccess(managerAndAdminRole, app.deleteLocationCheckInHandler),
	)
	mux.HandleFunc(
		"GET /locations/{locationId}",
		app.authAccess(allRoles, app.getLocationHandler),
	)
	mux.HandleFunc(
		"GET /all-locations",
		app.authAccess(managerAndAdminRole, app.getAllLocationsHandler),
	)
	mux.HandleFunc(
		"GET /locations",
		app.authAccess(managerAndAdminRole, app.getPaginatedLocationsHandler),
	)
	mux.HandleFunc(
		"POST /locations",
		app.authAccess(managerAndAdminRole, app.createLocationHandler),
	)
	mux.HandleFunc(
		"PATCH /locations/{locationId}",
		app.authAccess(allRoles, app.updateLocationHandler),
	)
	mux.HandleFunc(
		"DELETE /locations/{locationId}",
		app.authAccess(managerAndAdminRole, app.deleteLocationHandler),
	)
}

// @Summary	Get all check-ins at location for a specified day in a specified format
// @Tags		locations
// @Param		ids			query		[]string	true	"Location IDs"
// @Param		returnType	query		string		true	"ReturnType ('raw' or 'csv')"
// @Param		date		query		string		true	"Date (format: 'yyyy-MM-dd')"
// @Success	200			{object}	CheckInsGraphDto
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	404			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/all-locations/checkins/day [get].
func (app *Application) getLocationCheckInsDayHandler(w http.ResponseWriter,
	r *http.Request) {
	ids, err := parse.RequiredArrayQueryParam(r, "ids", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	returnType, err := parse.RequiredQueryParam[string](r, "returnType", nil)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	date, err := parse.RequiredQueryParam(
		r,
		"date",
		parse.Date(constants.DateFormat),
	)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	dateStrings, capacities, valueMap, err := app.services.Locations.GetCheckInsEntriesDay(
		r.Context(),
		user,
		ids,
		date,
	)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	if returnType == "csv" {
		filename := app.getTimeNowUTC().
			In(date.Location()).
			Format(constants.CSVFileNameFormat)
		filename = "Day-" + filename

		err = httptools.WriteCSV(
			w,
			filename,
			getCSVHeaders(valueMap),
			getCSVData(dateStrings, capacities, valueMap),
		)
	} else {
		err = httptools.WriteJSON(w, http.StatusOK, dtos.CheckInsGraphDto{
			Dates:                 dateStrings,
			CapacitiesPerLocation: capacities,
			ValuesPerSchool:       valueMap,
		}, nil)
	}

	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Get all check-ins at location for a specified range in a specified format
// @Tags		locations
// @Param		ids			query		[]string	true	"Location IDs"
// @Param		returnType	query		string		true	"ReturnType ('raw' or 'csv')"
// @Param		startDate	query		string		true	"StartDate (format: 'yyyy-MM-dd')"
// @Param		endDate		query		string		true	"EndDate (format: 'yyyy-MM-dd')"
// @Success	200			{object}	CheckInsGraphDto
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	404			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/all-locations/checkins/range [get].
func (app *Application) getLocationCheckInsRangeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ids, err := parse.RequiredArrayQueryParam(r, "ids", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	returnType, err := parse.RequiredQueryParam[string](r, "returnType", nil)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	startDate, err := parse.RequiredQueryParam(
		r,
		"startDate",
		parse.Date(constants.DateFormat),
	)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	endDate, err := parse.RequiredQueryParam(
		r,
		"endDate",
		parse.Date(constants.DateFormat),
	)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	//nolint:lll //it is what it is
	dateStrings, capacities, valueMap, err := app.services.Locations.GetCheckInsEntriesRange(
		r.Context(),
		user,
		ids,
		startDate,
		endDate,
	)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	if returnType == "csv" {
		filename := app.getTimeNowUTC().
			In(startDate.Location()).
			Format(constants.CSVFileNameFormat)
		filename = "Range-" + filename

		err = httptools.WriteCSV(
			w,
			filename,
			getCSVHeaders(valueMap),
			getCSVData(dateStrings, capacities, valueMap),
		)
	} else {
		err = httptools.WriteJSON(w, http.StatusOK, dtos.CheckInsGraphDto{
			Dates:                 dateStrings,
			CapacitiesPerLocation: capacities,
			ValuesPerSchool:       valueMap,
		}, nil)
	}

	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func getCSVHeaders(
	valueMap map[string][]int,
) []string {
	headers := []string{
		"datetime",
		"capacity",
	}

	for schoolName := range valueMap {
		headers = append(headers, schoolName)
	}

	return headers
}

func getCSVData(
	dateStrings []string,
	capacities map[string][]int,
	valuesPerSchool map[string][]int,
) [][]string {
	var output [][]string

	for i, dateString := range dateStrings {
		for _, values := range valuesPerSchool {
			var entry []string

			var totalCapacity int
			for _, capacity := range capacities {
				totalCapacity += capacity[i]
			}

			entry = append(entry, dateString)
			entry = append(entry, fmt.Sprintf("%d", totalCapacity))
			entry = append(entry, strconv.Itoa(values[i]))
			output = append(output, entry)
		}
	}

	return output
}

// @Summary	Get all checkins today
// @Tags		locations
// @Success	200	{object}	[]dtos.CheckInDto
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/locations/{id}/checkins [get].
func (app *Application) getAllCheckInsTodayHandler(w http.ResponseWriter,
	r *http.Request) {
	id, err := parse.URLParam(r, "locationId", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	_, checkIns, err := app.services.Locations.GetAllCheckInsOfDay(
		r.Context(),
		user,
		false,
		[]string{id},
		app.getTimeNowUTC(),
	)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, checkIns, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Delete check-in that occured today
// @Tags		locations
// @Param		locationId	path		string	true	"Location ID"
// @Param		checkInId	path		int		true	"Check-In ID"
// @Success	200			{object}	dtos.CheckInDto
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	404			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/locations/{locationId}/checkins/{checkInId} [delete].
func (app *Application) deleteLocationCheckInHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	locationID, err := parse.URLParam(r, "locationId", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	checkInID, err := parse.URLParam(r, "checkInId", parse.Int64(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	checkIn, err := app.services.Locations.DeleteCheckIn(
		r.Context(),
		user,
		locationID,
		checkInID,
	)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, checkIn, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Get single location
// @Tags		locations
// @Param		id	path		string	true	"Location ID"
// @Success	200	{object}	models.Location
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	404	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/locations/{id} [get].
func (app *Application) getLocationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam(r, "locationId", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	location, err := app.services.Locations.GetByID(r.Context(), user, id)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Get all locations paginated
// @Tags		locations
// @Param		page	query		int	false	"Page to fetch"
// @Success	200		{object}	dtos.PaginatedLocationsDto
// @Failure	400		{object}	ErrorDto
// @Failure	401		{object}	ErrorDto
// @Failure	500		{object}	ErrorDto
// @Router		/locations [get].
func (app *Application) getPaginatedLocationsHandler(w http.ResponseWriter,
	r *http.Request) {
	var pageSize int64 = 3

	page, err := parse.QueryParam(r, "page", 1, parse.Int64(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	result, err := getAllPaginated(
		r.Context(),
		app.services.Locations,
		user,
		page,
		pageSize,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Get all locations
// @Tags		locations
// @Success	200	{object}	[]models.Location
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/all-locations [get].
func (app *Application) getAllLocationsHandler(w http.ResponseWriter,
	r *http.Request) {
	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	locations, err := app.services.Locations.GetAll(r.Context(), user, false)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, locations, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Create location
// @Tags		locations
// @Param		createLocationDto	body		CreateLocationDto	true	"CreateLocationDto"
// @Success	201					{object}	models.Location
// @Failure	400					{object}	ErrorDto
// @Failure	401					{object}	ErrorDto
// @Failure	409					{object}	ErrorDto
// @Failure	500					{object}	ErrorDto
// @Router		/locations [post].
func (app *Application) createLocationHandler(w http.ResponseWriter, r *http.Request) {
	var createLocationDto dtos.CreateLocationDto

	err := httptools.ReadJSON(r.Body, &createLocationDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	if v, validationErrors := createLocationDto.Validate(); !v {
		httptools.FailedValidationResponse(w, r, validationErrors)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	location, err := app.services.Locations.Create(
		r.Context(),
		user,
		createLocationDto,
	)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusCreated, location, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Update location
// @Tags		locations
// @Param		id					path		string				true	"Location ID"
// @Param		updateLocationDto	body		UpdateLocationDto	true	"UpdateLocationDto"
// @Success	200					{object}	models.Location
// @Failure	400					{object}	ErrorDto
// @Failure	401					{object}	ErrorDto
// @Failure	409					{object}	ErrorDto
// @Failure	500					{object}	ErrorDto
// @Router		/locations/{id} [patch].
func (app *Application) updateLocationHandler(w http.ResponseWriter,
	r *http.Request) {
	var updateLocationDto dtos.UpdateLocationDto

	id, err := parse.URLParam(r, "locationId", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	err = httptools.ReadJSON(r.Body, &updateLocationDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	if v, validationErrors := updateLocationDto.Validate(); !v {
		httptools.FailedValidationResponse(w, r, validationErrors)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	location, err := app.services.Locations.Update(
		r.Context(),
		user,
		id,
		updateLocationDto,
	)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Delete location
// @Tags		locations
// @Param		id	path		string	true	"Location ID"
// @Success	200	{object}	models.Location
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	404	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/locations/{id} [delete].
func (app *Application) deleteLocationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam(r, "locationId", parse.UUID)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	location, err := app.services.Locations.Delete(r.Context(), user, id)
	if err != nil {
		httptools.HandleError(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
