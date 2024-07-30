package main

import (
	"errors"
	"net/http"
	"time"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	"github.com/xdoubleu/essentia/pkg/context"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/parse"
	"github.com/xdoubleu/essentia/pkg/tools"

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
// @Success	200			{object}	[]CheckInsLocationEntryRaw
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
		parse.DateFunc(constants.DateFormat),
	)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	startDate := tools.StartOfDay(date)
	endDate := tools.EndOfDay(date)

	user := context.GetContextValue[models.User](r.Context(), userContextKey)

	for _, id := range ids {
		var location *models.Location
		location, err = app.services.Locations.GetByID(r.Context(), id)
		if err != nil ||
			(user.Role == models.DefaultRole && location.UserID != user.ID) {
			httptools.NotFoundResponse(w, r, err, "location", id, "id")
			return
		}
	}

	checkIns, err := app.services.CheckIns.GetAllInRange(
		r.Context(),
		ids,
		startDate,
		endDate,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	checkInEntries := app.services.Locations.GetCheckInsEntriesDay(
		checkIns,
		schools,
	)

	if returnType == "csv" {
		filename := time.Now().
			In(startDate.Location()).
			Format(constants.CSVFileNameFormat)
		filename = "Day-" + filename

		data := dtos.ConvertCheckInsLocationEntryRawMapToCSV(
			checkInEntries,
		)
		err = httptools.WriteCSV(w, filename, data)
	} else {
		err = httptools.WriteJSON(w, http.StatusOK, checkInEntries, nil)
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
// @Success	200			{object}	[]CheckInsLocationEntryRaw
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
		parse.DateFunc(constants.DateFormat),
	)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	endDate, err := parse.RequiredQueryParam(
		r,
		"endDate",
		parse.DateFunc(constants.DateFormat),
	)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	startDate = tools.StartOfDay(startDate)
	endDate = tools.EndOfDay(endDate)

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	user := context.GetContextValue[models.User](r.Context(), userContextKey)

	for _, id := range ids {
		var location *models.Location
		location, err = app.services.Locations.GetByID(r.Context(), id)
		if err != nil ||
			(user.Role == models.DefaultRole && location.UserID != user.ID) {
			httptools.NotFoundResponse(w, r, err, "location", id, "id")
			return
		}
	}

	checkIns, err := app.services.CheckIns.GetAllInRange(
		r.Context(),
		ids,
		startDate,
		endDate,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	checkInEntries := app.services.Locations.GetCheckInsEntriesRange(
		startDate,
		endDate,
		checkIns,
		schools,
	)

	if returnType == "csv" {
		filename := time.Now().
			In(startDate.Location()).
			Format(constants.CSVFileNameFormat)
		filename = "Range-" + filename

		data := dtos.ConvertCheckInsLocationEntryRawMapToCSV(
			checkInEntries,
		)
		err = httptools.WriteCSV(w, filename, data)
	} else {
		err = httptools.WriteJSON(w, http.StatusOK, checkInEntries, nil)
	}

	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
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

	user := context.GetContextValue[models.User](r.Context(), userContextKey)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		httptools.NotFoundResponse(w, r, err, "location", id, "id")
		return
	}

	loc, _ := time.LoadLocation(location.TimeZone)
	today := time.Now().In(loc)
	startOfToday := tools.StartOfDay(today)
	endOfToday := tools.EndOfDay(today)

	checkIns, err := app.services.CheckIns.GetAllInRange(
		r.Context(),
		[]string{location.ID},
		startOfToday,
		endOfToday,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	schoolMap, _ := app.services.Schools.GetSchoolMaps(schools)

	checkInDtos := make([]dtos.CheckInDto, 0)
	for _, checkIn := range checkIns {
		checkInDto := dtos.CheckInDto{
			ID:         checkIn.ID,
			LocationID: checkIn.LocationID,
			SchoolName: schoolMap[checkIn.SchoolID],
			Capacity:   checkIn.Capacity,
			CreatedAt:  checkIn.CreatedAt,
		}
		checkInDtos = append(checkInDtos, checkInDto)
	}

	err = httptools.WriteJSON(w, http.StatusOK, checkInDtos, nil)
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

	checkInID, err := parse.URLParam(r, "checkInId", parse.Int64Func(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	location, err := app.services.Locations.GetByID(r.Context(), locationID)
	if err != nil {
		httptools.NotFoundResponse(w, r, err, "location", locationID, "id")
		return
	}

	checkIn, err := app.services.CheckIns.GetByID(r.Context(), location, checkInID)
	if err != nil {
		httptools.NotFoundResponse(w, r, err, "checkIn", checkInID, "id")
		return
	}

	today := tools.TimeZoneIndependentTimeNow(location.TimeZone)
	startOfToday := tools.StartOfDay(today)
	endOfToday := tools.EndOfDay(today)

	if !(checkIn.CreatedAt.Time.After(startOfToday) &&
		checkIn.CreatedAt.Time.Before(endOfToday)) {
		httptools.BadRequestResponse(
			w,
			r,
			errors.New("checkIn didn't occur today and thus can't be deleted"),
		)
		return
	}

	err = app.services.CheckIns.Delete(r.Context(), checkIn.ID)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}
	schoolMap, _ := app.services.Schools.GetSchoolMaps(schools)

	checkInDto := dtos.CheckInDto{
		ID:         checkIn.ID,
		LocationID: checkIn.LocationID,
		SchoolName: schoolMap[checkIn.SchoolID],
		Capacity:   checkIn.Capacity,
		CreatedAt:  checkIn.CreatedAt,
	}

	err = httptools.WriteJSON(w, http.StatusOK, checkInDto, nil)
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

	user := context.GetContextValue[models.User](r.Context(), userContextKey)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		httptools.NotFoundResponse(w, r, err, "location", id, "id")
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

	page, err := parse.QueryParam(r, "page", 1, parse.Int64Func(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	result, err := getAllPaginated(
		r.Context(),
		app.services.Locations,
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
	locations, err := app.services.Locations.GetAll(r.Context())
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

	if v := createLocationDto.Validate(); !v.Valid() {
		httptools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	existingLocation, err := app.services.Locations.GetByName(
		r.Context(),
		createLocationDto.Name,
	)
	if existingLocation != nil || !errors.Is(err, errortools.ErrResourceNotFound) {
		httptools.ConflictResponse(
			w,
			r,
			err,
			"location",
			createLocationDto.Name,
			"name",
		)
		return
	}

	existingUser, err := app.services.Users.GetByUsername(
		r.Context(),
		createLocationDto.Username,
	)
	if existingUser != nil || !errors.Is(err, errortools.ErrResourceNotFound) {
		httptools.ConflictResponse(
			w,
			r,
			err,
			"user",
			createLocationDto.Username,
			"username",
		)
		return
	}

	location, err := app.services.Locations.Create(
		r.Context(),
		createLocationDto.Name,
		createLocationDto.Capacity,
		createLocationDto.TimeZone,
		createLocationDto.Username,
		createLocationDto.Password,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
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

	if v := updateLocationDto.Validate(); !v.Valid() {
		httptools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	hasConflicts := app.checkForConflictsOnUpdate(w, r, updateLocationDto)
	if hasConflicts {
		return
	}

	user := context.GetContextValue[models.User](r.Context(), userContextKey)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		httptools.NotFoundResponse(w, r, err, "location", id, "id")
		return
	}

	locationUser, err := app.services.Users.GetByID(
		r.Context(),
		location.UserID,
		models.DefaultRole,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = app.services.Locations.Update(
		r.Context(),
		location,
		locationUser,
		updateLocationDto,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

func (app *Application) checkForConflictsOnUpdate(
	w http.ResponseWriter,
	r *http.Request,
	updateLocationDto dtos.UpdateLocationDto,
) bool {
	if updateLocationDto.Name != nil {
		existingLocation, err := app.services.Locations.GetByName(
			r.Context(),
			*updateLocationDto.Name,
		)

		if existingLocation != nil || !errors.Is(err, errortools.ErrResourceNotFound) {
			httptools.ConflictResponse(
				w,
				r,
				err,
				"location",
				*updateLocationDto.Name,
				"name",
			)
			return true
		}
	}

	if updateLocationDto.Username != nil {
		existingUser, err := app.services.Users.GetByUsername(
			r.Context(),
			*updateLocationDto.Username,
		)

		if existingUser != nil || !errors.Is(err, errortools.ErrResourceNotFound) {
			httptools.ConflictResponse(
				w,
				r,
				err,
				"user",
				*updateLocationDto.Username,
				"username",
			)
			return true
		}
	}

	return false
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

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil {
		httptools.NotFoundResponse(w, r, err, "location", id, "id")
		return
	}

	err = app.services.Locations.Delete(r.Context(), location)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
