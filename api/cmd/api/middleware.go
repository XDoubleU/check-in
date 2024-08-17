package main

import (
	"errors"
	"net/http"

	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/models"
)

func (app *Application) authAccess(allowedRoles []models.Role,
	next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("accessToken")

		if err != nil {
			httptools.UnauthorizedResponse(w, r,
				errortools.NewUnauthorizedError(errors.New("no token in cookies")))
			return
		}

		_, user, err := app.services.Auth.GetToken(
			r.Context(),
			models.AccessScope,
			tokenCookie.Value,
		)
		if err != nil {
			httptools.HandleError(w, r, err, nil)
			return
		}

		if user.Role == models.DefaultRole {
			user, err = app.services.Locations.GetDefaultUserByUserID(
				r.Context(),
				user.ID,
			)
			if err != nil {
				httptools.ServerErrorResponse(w, r, err)
			}
		}

		r = r.WithContext(app.contextSetUser(r.Context(), *user))

		forbidden := true
		for _, role := range allowedRoles {
			if user.Role == role {
				forbidden = false
				break
			}
		}

		if forbidden {
			httptools.ForbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) authRefresh(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("refreshToken")

		if err != nil {
			httptools.UnauthorizedResponse(w, r,
				errortools.NewUnauthorizedError(errors.New("no token in cookies")))
			return
		}

		token, user, err := app.services.Auth.GetToken(r.Context(),
			models.RefreshScope, tokenCookie.Value)
		if err != nil {
			httptools.HandleError(w, r, err, nil)
			return
		}

		r = r.WithContext(app.contextSetUser(r.Context(), *user))

		if token.Used {
			err = app.services.Auth.DeleteAllTokensForUser(r.Context(), user.ID)
			if err != nil {
				httptools.ServerErrorResponse(w, r, err)
			}
			httptools.UnauthorizedResponse(
				w,
				r,
				errortools.NewUnauthorizedError(errors.New("invalid token")),
			)
			return
		}

		err = app.services.Auth.SetTokenAsUsed(r.Context(), tokenCookie.Value)
		if err != nil {
			httptools.ServerErrorResponse(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
