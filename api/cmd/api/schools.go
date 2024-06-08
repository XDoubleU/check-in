package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/http_tools"
	"github.com/julienschmidt/httprouter"

	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
)

func (app *application) schoolsRoutes(router *httprouter.Router) {
	router.HandlerFunc(
		http.MethodGet,
		"/schools",
		app.authAccess(managerAndAdminRole, app.getPaginatedSchoolsHandler),
	)
	router.HandlerFunc(
		http.MethodPost,
		"/schools",
		app.authAccess(managerAndAdminRole, app.createSchoolHandler),
	)
	router.HandlerFunc(
		http.MethodPatch,
		"/schools/:id",
		app.authAccess(managerAndAdminRole, app.updateSchoolHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/schools/:id",
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

	page, err := helpers.ReadIntQueryParam(r, "page", 1)
	if err != nil {
		http_tools.BadRequestResponse(w, r, err)
		return
	}

	result, err := getAllPaginated[models.School](
		r.Context(),
		app.services.Schools,
		page,
		pageSize,
	)
	if err != nil {
		http_tools.ServerErrorResponse(w, r, err)
		return
	}

	err = http_tools.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		http_tools.ServerErrorResponse(w, r, err)
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

	err := http_tools.ReadJSON(r.Body, &schoolDto)
	if err != nil {
		http_tools.BadRequestResponse(w, r, err)
		return
	}

	if v := schoolDto.Validate(); !v.Valid() {
		http_tools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	school, err := app.services.Schools.Create(r.Context(), schoolDto.Name)
	if err != nil {
		http_tools.ConflictResponse(w, r, err, "school", schoolDto.Name, "name")
		return
	}

	err = http_tools.WriteJSON(w, http.StatusCreated, school, nil)
	if err != nil {
		http_tools.ServerErrorResponse(w, r, err)
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

	id, err := helpers.ReadIntURLParam(r, "id")
	if err != nil {
		http_tools.BadRequestResponse(w, r, err)
		return
	}

	err = http_tools.ReadJSON(r.Body, &schoolDto)
	if err != nil {
		http_tools.BadRequestResponse(w, r, err)
		return
	}

	if v := schoolDto.Validate(); !v.Valid() {
		http_tools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	school, err := app.services.Schools.GetByIDWithoutReadOnly(r.Context(), id)
	if err != nil {
		http_tools.NotFoundResponse(w, r, err, "school", id, "id")
		return
	}

	err = app.services.Schools.Update(r.Context(), school, schoolDto)
	if err != nil {
		http_tools.ConflictResponse(w, r, err, "school", schoolDto.Name, "name")
		return
	}

	err = http_tools.WriteJSON(w, http.StatusOK, school, nil)
	if err != nil {
		http_tools.ServerErrorResponse(w, r, err)
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
	id, err := helpers.ReadIntURLParam(r, "id")
	if err != nil {
		http_tools.BadRequestResponse(w, r, err)
		return
	}

	school, err := app.services.Schools.GetByIDWithoutReadOnly(r.Context(), id)
	if err != nil {
		http_tools.NotFoundResponse(w, r, err, "school", id, "id")
		return
	}

	err = app.services.Schools.Delete(r.Context(), school.ID)
	if err != nil {
		http_tools.ServerErrorResponse(w, r, err)
		return
	}

	err = http_tools.WriteJSON(w, http.StatusOK, school, nil)
	if err != nil {
		http_tools.ServerErrorResponse(w, r, err)
	}
}
