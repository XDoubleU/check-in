package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"check-in/api/internal/dtos"
	"check-in/api/internal/helpers"
	"check-in/api/internal/models"
	"check-in/api/internal/validator"
)

func (app *application) usersRoutes(router *httprouter.Router) {
	router.HandlerFunc(
		http.MethodGet,
		"/current-user",
		app.authAccess(allRoles, app.getInfoLoggedInUserHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/users",
		app.authAccess(adminRole, app.getPaginatedManagerUsersHandler),
	)
	router.HandlerFunc(
		http.MethodGet,
		"/users/:id",
		app.authAccess(managerAndAdminRole, app.getUserHandler),
	)
	router.HandlerFunc(
		http.MethodPost,
		"/users",
		app.authAccess(adminRole, app.createManagerUserHandler),
	)
	router.HandlerFunc(
		http.MethodPatch,
		"/users/:id",
		app.authAccess(adminRole, app.updateManagerUserHandler),
	)
	router.HandlerFunc(
		http.MethodDelete,
		"/users/:id",
		app.authAccess(adminRole, app.deleteManagerUserHandler),
	)
}

// @Summary	Get info of logged in user
// @Tags		users
// @Success	200	{object}	User
// @Failure	401	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/current-user [get].
func (app *application) getInfoLoggedInUserHandler(w http.ResponseWriter,
	r *http.Request) {
	user := app.contextGetUser(r)

	err := helpers.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary	Get single user
// @Tags		users
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	User
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	404	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/users/{id} [get].
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadUUIDURLParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.services.Users.GetByID(r.Context(), id, models.DefaultRole)
	if err != nil {
		app.notFoundResponse(w, r, err, "user", "id", id, "id")
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary	Get all users paginated
// @Tags		users
// @Param		page	query		int	false	"Page to fetch"
// @Success	200		{object}	PaginatedUsersDto
// @Failure	400		{object}	ErrorDto
// @Failure	401		{object}	ErrorDto
// @Failure	500		{object}	ErrorDto
// @Router		/users [get].
func (app *application) getPaginatedManagerUsersHandler(w http.ResponseWriter,
	r *http.Request) {
	var pageSize int64 = 4

	page, err := helpers.ReadIntQueryParam(r, "page", 1)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	result, err := getAllPaginated[models.User](
		r.Context(),
		app.services.Users,
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

// @Summary	Create user
// @Tags		users
// @Param		createUserDto	body		CreateUserDto	true	"CreateUserDto"
// @Success	201				{object}	User
// @Failure	400				{object}	ErrorDto
// @Failure	401				{object}	ErrorDto
// @Failure	409				{object}	ErrorDto
// @Failure	500				{object}	ErrorDto
// @Router		/users [post].
func (app *application) createManagerUserHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var createUserDto dtos.CreateUserDto

	err := helpers.ReadJSON(r.Body, &createUserDto)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if dtos.ValidateCreateUserDto(v, createUserDto); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.services.Users.Create(
		r.Context(),
		createUserDto.Username,
		createUserDto.Password,
		models.ManagerRole,
	)
	if err != nil {
		app.conflictResponse(
			w,
			r,
			err,
			"user",
			"username",
			createUserDto.Username,
			"username",
		)
		return
	}

	err = helpers.WriteJSON(w, http.StatusCreated, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary	Update user
// @Tags		users
// @Param		id				path		string			true	"User ID"
// @Param		updateUserDto	body		UpdateUserDto	true	"UpdateUserDto"
// @Success	200				{object}	User
// @Failure	400				{object}	ErrorDto
// @Failure	401				{object}	ErrorDto
// @Failure	409				{object}	ErrorDto
// @Failure	500				{object}	ErrorDto
// @Router		/users/{id} [patch].
func (app *application) updateManagerUserHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var updateUserDto dtos.UpdateUserDto

	id, err := helpers.ReadUUIDURLParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = helpers.ReadJSON(r.Body, &updateUserDto)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if dtos.ValidateUpdateUserDto(v, updateUserDto); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.services.Users.GetByID(r.Context(), id, models.ManagerRole)
	if err != nil {
		app.notFoundResponse(w, r, err, "user", "id", id, "id")
		return
	}

	err = app.services.Users.Update(
		r.Context(),
		user,
		updateUserDto,
		models.ManagerRole,
	)
	if err != nil {
		app.conflictResponse(
			w,
			r,
			err,
			"user",
			"username",
			*updateUserDto.Username,
			"username",
		)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary	Delete user
// @Tags		users
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	User
// @Failure	400	{object}	ErrorDto
// @Failure	401	{object}	ErrorDto
// @Failure	404	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/users/{id} [delete].
func (app *application) deleteManagerUserHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := helpers.ReadUUIDURLParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.services.Users.GetByID(r.Context(), id, models.ManagerRole)
	if err != nil {
		app.notFoundResponse(w, r, err, "user", "id", id, "id")
		return
	}

	err = app.services.Users.Delete(r.Context(), user.ID, models.ManagerRole)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
