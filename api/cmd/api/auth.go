package main

import (
	"net/http"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	"github.com/xdoubleu/essentia/pkg/config"
	contexttools "github.com/xdoubleu/essentia/pkg/context"

	"check-in/api/internal/constants"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func (app *Application) authRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/signin", app.signInHandler)
	mux.HandleFunc(
		"GET /auth/signout",
		app.authAccess(allRoles, app.signOutHandler),
	)
	mux.HandleFunc(
		"GET /auth/refresh",
		app.authRefresh(app.refreshHandler),
	)
}

// @Summary	Sign in a user
// @Tags		auth
// @Param		signInDto	body		SignInDto	true	"SignInDto"
// @Success	200			{object}	User
// @Failure	400			{object}	ErrorDto
// @Failure	401			{object}	ErrorDto
// @Failure	500			{object}	ErrorDto
// @Router		/auth/signin [post].
func (app *Application) signInHandler(w http.ResponseWriter, r *http.Request) {
	var signInDto *dtos.SignInDto

	err := httptools.ReadJSON(r.Body, &signInDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	user, err := app.services.Auth.SignInUser(r.Context(), signInDto)
	if err != nil {
		httptools.HandleError(w, r, err, signInDto.ValidationErrors)
		return
	}

	secure := app.config.Env == config.ProdEnv
	accessTokenCookie, err := app.services.Auth.CreateCookie(
		r.Context(),
		models.AccessScope,
		user.ID,
		app.config.AccessExpiry,
		secure,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, accessTokenCookie)

	if user.Role != models.AdminRole && signInDto.RememberMe {
		var refreshTokenCookie *http.Cookie
		refreshTokenCookie, err = app.services.Auth.CreateCookie(
			r.Context(),
			models.RefreshScope,
			user.ID,
			app.config.RefreshExpiry,
			secure,
		)
		if err != nil {
			httptools.ServerErrorResponse(w, r, err)
			return
		}

		http.SetCookie(w, refreshTokenCookie)
	}

	err = httptools.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}

// @Summary	Sign out a user
// @Tags		auth
// @Success	200	{object}	nil
// @Failure	401	{object}	ErrorDto
// @Router		/auth/signout [get].
func (app *Application) signOutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, _ := r.Cookie("accessToken")
	refreshToken, _ := r.Cookie("refreshToken")

	deleteAccessToken, err := app.services.Auth.DeleteCookie(
		models.AccessScope,
		accessToken.Value,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, deleteAccessToken)

	if refreshToken == nil {
		return
	}

	deleteRefreshToken, err := app.services.Auth.DeleteCookie(
		models.RefreshScope,
		refreshToken.Value,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, deleteRefreshToken)
}

// @Summary	Refresh access token
// @Tags		auth
// @Success	200	{object}	nil
// @Failure	401	{object}	ErrorDto
// @Failure	500	{object}	ErrorDto
// @Router		/auth/refresh [get].
func (app *Application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	user := contexttools.GetValue[models.User](r.Context(), constants.UserContextKey)
	secure := app.config.Env == config.ProdEnv

	accessTokenCookie, err := app.services.Auth.CreateCookie(
		r.Context(),
		models.AccessScope,
		user.ID,
		app.config.AccessExpiry,
		secure,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, accessTokenCookie)

	refreshTokenCookie, err := app.services.Auth.CreateCookie(
		r.Context(),
		models.RefreshScope,
		user.ID,
		app.config.RefreshExpiry,
		secure,
	)
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, refreshTokenCookie)

	err = app.services.Auth.DeleteExpiredTokens(r.Context())
	if err != nil {
		httptools.ServerErrorResponse(w, r, err)
	}
}
