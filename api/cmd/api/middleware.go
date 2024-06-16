package main

import (
	"errors"
	"net/http"

	"check-in/api/internal/models"

	"github.com/XDoubleU/essentia/pkg/http_tools"
)

func (app *application) authAccess(allowedRoles []models.Role,
	next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("accessToken")

		if err != nil {
			http_tools.UnauthorizedResponse(w, r, "No token in cookies")
			return
		}

		_, user, err := app.repositories.Auth.GetToken(
			r.Context(),
			models.AccessScope,
			tokenCookie.Value,
		)
		if err != nil {
			switch {
			case errors.Is(err, http_tools.ErrRecordNotFound):
				http_tools.UnauthorizedResponse(w, r, "Invalid token")
			default:
				http_tools.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		forbidden := true
		for _, role := range allowedRoles {
			if user.Role == role {
				forbidden = false
				break
			}
		}

		if forbidden {
			http_tools.ForbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authRefresh(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("refreshToken")

		if err != nil {
			http_tools.UnauthorizedResponse(w, r, "No token in cookies")
			return
		}

		token, user, err := app.repositories.Auth.GetToken(r.Context(),
			models.RefreshScope, tokenCookie.Value)
		if err != nil {
			switch {
			case errors.Is(err, http_tools.ErrRecordNotFound):
				http_tools.UnauthorizedResponse(w, r, "Invalid token")
			default:
				http_tools.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		if token.Used {
			err = app.repositories.Auth.DeleteAllTokensForUser(r.Context(), user.ID)
			if err != nil {
				panic(err)
			}
			http_tools.UnauthorizedResponse(w, r, "Invalid token")
			return
		}

		err = app.repositories.Auth.SetTokenAsUsed(r.Context(), tokenCookie.Value)
		if err != nil {
			http_tools.ServerErrorResponse(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
