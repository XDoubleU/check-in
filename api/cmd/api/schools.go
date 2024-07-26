package main

import (
	"net/http"

	"github.com/xdoubleu/essentia/pkg/httptools"
	"github.com/xdoubleu/essentia/pkg/parse"

	"check-in/api/internal/dtos"
)

func (app *application) schoolsRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /schools",
		app.authAccess(managerAndAdminRole, app.getPaginatedSchoolsHandler),
	)
	mux.HandleFunc(
		"POST /schools",
		app.authAccess(managerAndAdminRole, app.createSchoolHandler),
	)
	mux.HandleFunc(
		"PATCH /schools/:id",
		app.authAccess(managerAndAdminRole, app.updateSchoolHandler),
	)
	mux.HandleFunc(
		"DELETE /schools/:id",
		app.authAccess(managerAndAdminRole, app.deleteSchoolHandler),
	)
}

// @Summary	Get all schools paginated
// @Tags		schools
// @Param		page	query		int	false	"Page to fetch"
// @Success	200		{object}	PaginatedSchoolsDto
// @Failure	400		{object}	ErrorDto
// @Failure	401		{object}	ErrorDto
// @Failure	500		{object}	ErrorDto
// @Router		/schools [get].
func (app *application) getPaginatedSchoolsHandler(w http.ResponseWriter,
	r *http.Request) {
	var pageSize int64 = 4

	page, err := parse.QueryParam(r, "page", 1, parse.Int64Func(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	result, err := getAllPaginated(
		r.Context(),
		app.services.Schools,
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

// @Summary	Create school
// @Tags		schools
// @Param		schoolDto	body		SchoolDto	true	"SchoolDto"
// @Success	201			{object}	School
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	409			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/schools [post].
func (app *application) createSchoolHandler(w http.ResponseWriter, r *http.Request) {
	var schoolDto dtos.SchoolDto

	err := httptools.ReadJSON(r.Body, &schoolDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	if v := schoolDto.Validate(); !v.Valid() {
		httptools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	school, err := app.services.Schools.Create(r.Context(), schoolDto.Name)
	if err != nil {
		httptools.ConflictResponse(w, r, err, "school", schoolDto.Name, "name")
		return
	}

	err = httptools.WriteJSON(w, http.StatusCreated, school, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Update school
// @Tags		schools
// @Param		id			path		int			true	"School ID"
// @Param		schoolDto	body		SchoolDto	true	"SchoolDto"
// @Success	200			{object}	School
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	409			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/schools/{id} [patch].
func (app *application) updateSchoolHandler(w http.ResponseWriter, r *http.Request) {
	var schoolDto dtos.SchoolDto

	id, err := parse.URLParam(r, "id", parse.Int64Func(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	err = httptools.ReadJSON(r.Body, &schoolDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	if v := schoolDto.Validate(); !v.Valid() {
		httptools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	school, err := app.services.Schools.GetByIDWithoutReadOnly(r.Context(), id)
	if err != nil {
		httptools.NotFoundResponse(w, r, err, "school", id, "id")
		return
	}

	err = app.services.Schools.Update(r.Context(), school, schoolDto)
	if err != nil {
		httptools.ConflictResponse(w, r, err, "school", schoolDto.Name, "name")
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, school, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Delete school
// @Tags		schools
// @Param		id	path		int	true	"School ID"
// @Success	200	{object}	School
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	404	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/schools/{id} [delete].
func (app *application) deleteSchoolHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam(r, "id", parse.Int64Func(true, false))
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	school, err := app.services.Schools.GetByIDWithoutReadOnly(r.Context(), id)
	if err != nil {
		httptools.NotFoundResponse(w, r, err, "school", id, "id")
		return
	}

	err = app.services.Schools.Delete(r.Context(), school.ID)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	err = httptools.WriteJSON(w, http.StatusOK, school, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
