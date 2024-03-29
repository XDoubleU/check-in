package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/services"
	"check-in/api/internal/validator"
)

func (app *application) locationsRoutes(router *httprouter.Router) {
	router.HandlerFunc(
		http.MethodGet,
		"/all-locations/checkins/range",
		app.authAccess(allRoles, app.getLocationCheckInsRangeHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/all-locations/checkins/day",
		app.authAccess(allRoles, app.getLocationCheckInsDayHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/locations/:locationId/checkins",
		app.authAccess(allRoles, app.getAllCheckInsTodayHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/locations/:locationId/checkins/:checkInId",
		app.authAccess(managerAndAdminRole, app.deleteLocationCheckInHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/locations/:locationId",
		app.authAccess(allRoles, app.getLocationHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/all-locations",
		app.authAccess(managerAndAdminRole, app.getAllLocationsHandler),
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
		"/locations/:locationId",
		app.authAccess(allRoles, app.updateLocationHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/locations/:locationId",
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
func (app *application) getLocationCheckInsDayHandler(w http.ResponseWriter,
	r *http.Request) {
	ids, err := helpers.ReadUUIDArrayQueryParam(r, "ids")
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

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	startDate := helpers.StartOfDay(date)
	endDate := helpers.EndOfDay(date)

	user := app.contextGetUser(r)

	for _, id := range ids {
		var location *models.Location
		location, err = app.services.Locations.GetByID(r.Context(), id)
		if err != nil ||
			(user.Role == models.DefaultRole && location.UserID != user.ID) {
			app.notFoundResponse(w, r, err, "location", "id", id, "id")
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
		app.serverErrorResponse(w, r, err)
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
		err = helpers.WriteCSV(w, filename, data)
	} else {
		err = helpers.WriteJSON(w, http.StatusOK, checkInEntries, nil)
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
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
func (app *application) getLocationCheckInsRangeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ids, err := helpers.ReadUUIDArrayQueryParam(r, "ids")
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

	startDate = helpers.StartOfDay(startDate)
	endDate = helpers.EndOfDay(endDate)

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	for _, id := range ids {
		var location *models.Location
		location, err = app.services.Locations.GetByID(r.Context(), id)
		if err != nil ||
			(user.Role == models.DefaultRole && location.UserID != user.ID) {
			app.notFoundResponse(w, r, err, "location", "id", id, "id")
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
		filename := time.Now().
			In(startDate.Location()).
			Format(constants.CSVFileNameFormat)
		filename = "Range-" + filename

		data := dtos.ConvertCheckInsLocationEntryRawMapToCSV(
			checkInEntries,
		)
		err = helpers.WriteCSV(w, filename, data)
	} else {
		err = helpers.WriteJSON(w, http.StatusOK, checkInEntries, nil)
	}

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary	Get all checkins today
// @Tags		locations
// @Success	200	{object}	[]dtos.CheckInDto
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/locations/{id}/checkins [get].
func (app *application) getAllCheckInsTodayHandler(w http.ResponseWriter,
	r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r, "locationId")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		app.notFoundResponse(w, r, err, "location", "id", id, "id")
		return
	}

	loc, _ := time.LoadLocation(location.TimeZone)
	today := time.Now().In(loc)
	startOfToday := helpers.StartOfDay(&today)
	endOfToday := helpers.EndOfDay(&today)

	checkIns, err := app.services.CheckIns.GetAllInRange(
		r.Context(),
		[]string{location.ID},
		startOfToday,
		endOfToday,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

	err = helpers.WriteJSON(w, http.StatusOK, checkInDtos, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
func (app *application) deleteLocationCheckInHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	locationID, err := helpers.ReadUUIDURLParam(r, "locationId")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	checkInID, err := helpers.ReadIntURLParam(r, "checkInId")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	location, err := app.services.Locations.GetByID(r.Context(), locationID)
	if err != nil {
		app.notFoundResponse(w, r, err, "location", "id", locationID, "id")
		return
	}

	checkIn, err := app.services.CheckIns.GetByID(r.Context(), location, checkInID)
	if err != nil {
		app.notFoundResponse(w, r, err, "checkIn", "id", checkInID, "id")
		return
	}

	today := helpers.TimeZoneIndependentTimeNow(location.TimeZone)
	startOfToday := helpers.StartOfDay(&today)
	endOfToday := helpers.EndOfDay(&today)

	if !(checkIn.CreatedAt.Time.After(*startOfToday) &&
		checkIn.CreatedAt.Time.Before(*endOfToday)) {
		app.badRequestResponse(
			w,
			r,
			errors.New("checkIn didn't occur today and thus can't be deleted"),
		)
		return
	}

	err = app.services.CheckIns.Delete(r.Context(), checkIn.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schools, err := app.services.Schools.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

	err = helpers.WriteJSON(w, http.StatusOK, checkInDto, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
func (app *application) getLocationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r, "locationId")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil || (user.Role == models.DefaultRole && location.UserID != user.ID) {
		app.notFoundResponse(w, r, err, "location", "id", id, "id")
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

// @Summary	Get all locations
// @Tags		locations
// @Success	200	{object}	[]models.Location
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/all-locations [get].
func (app *application) getAllLocationsHandler(w http.ResponseWriter,
	r *http.Request) {
	locations, err := app.services.Locations.GetAll(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, locations, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
		app.conflictResponse(
			w,
			r,
			err,
			"location",
			"name",
			createLocationDto.Name,
			"name",
		)
		return
	}

	existingUser, err := app.services.Users.GetByUsername(
		r.Context(),
		createLocationDto.Username,
	)
	if existingUser != nil || !errors.Is(err, services.ErrRecordNotFound) {
		app.conflictResponse(
			w,
			r,
			err,
			"user",
			"username",
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
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusCreated, location, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
func (app *application) updateLocationHandler(w http.ResponseWriter,
	r *http.Request) {
	var updateLocationDto dtos.UpdateLocationDto

	id, err := helpers.ReadUUIDURLParam(r, "locationId")
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
		app.notFoundResponse(w, r, err, "location", "id", id, "id")
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
			app.conflictResponse(
				w,
				r,
				err,
				"location",
				"name",
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

		if existingUser != nil || !errors.Is(err, services.ErrRecordNotFound) {
			app.conflictResponse(
				w,
				r,
				err,
				"user",
				"username",
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
func (app *application) deleteLocationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r, "locationId")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	location, err := app.services.Locations.GetByID(r.Context(), id)
	if err != nil {
		app.notFoundResponse(w, r, err, "location", "id", id, "id")
		return
	}

	err = app.services.Locations.Delete(r.Context(), location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, location, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
