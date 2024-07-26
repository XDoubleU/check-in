package main

import (
	"errors"
	"net/http"

	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/contexttools"
	"github.com/xdoubleu/essentia/pkg/httptools"
	"github.com/xdoubleu/essentia/pkg/sentrytools"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

func (app *application) authRoutes(mux *http.ServeMux) {
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
func (app *application) signInHandler(w http.ResponseWriter, r *http.Request) {
	var signInDto dtos.SignInDto

	err := httptools.ReadJSON(r.Body, &signInDto)
	if err != nil {
		httptools.BadRequestResponse(w, r, err)
		return
	}

	if v := signInDto.Validate(); !v.Valid() {
		httptools.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.services.Users.GetByUsername(r.Context(), signInDto.Username)
	if err != nil {
		if errors.Is(err, httptools.ErrResourceNotFound) {
			httptools.UnauthorizedResponse(w, r, "Invalid Credentials")
		} else {
			httptools.ServerErrorResponse(w, r, err)
		}

		return
	}

	match, _ := user.CompareHashAndPassword(signInDto.Password)
	if !match {
		httptools.UnauthorizedResponse(w, r, "Invalid Credentials")
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
func (app *application) signOutHandler(w http.ResponseWriter, r *http.Request) {
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
func (app *application) refreshHandler(w http.ResponseWriter, r *http.Request) {
	user := contexttools.GetContextValue[models.User](r.Context(), userContextKey)
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

	go func() {
		sentrytools.GoRoutineErrorHandler(
			"delete expired tokens",
			app.services.Auth.DeleteExpiredTokens,
		)
	}()
}
