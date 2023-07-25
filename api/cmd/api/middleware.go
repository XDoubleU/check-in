package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"golang.org/x/time/rate"

	"check-in/api/internal/models"
	"check-in/api/internal/services"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", app.config.WebURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	var rps rate.Limit = 10
	var bucketSize = 30

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(rps, bucketSize)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authAccess(allowedRoles []models.Roles,
	next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("accessToken")

		if err != nil {
			app.unauthorizedResponse(w, r, "No token in cookies")
			return
		}

		_, user, err := app.services.Auth.GetToken(r.Context(),
			models.AccessScope, tokenCookie.Value)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrRecordNotFound):
				app.unauthorizedResponse(w, r, "Invalid token")
			default:
				app.serverErrorResponse(w, r, err)
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
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authRefresh(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("refreshToken")

		if err != nil {
			app.unauthorizedResponse(w, r, "No token in cookies")
			return
		}

		token, user, err := app.services.Auth.GetToken(r.Context(),
			models.RefreshScope, tokenCookie.Value)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrRecordNotFound):
				app.unauthorizedResponse(w, r, "Invalid token")
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		if token.Used {
			err = app.services.Auth.DeleteAllTokensForUser(r.Context(), user.ID)
			if err != nil {
				panic(err)
			}
			app.unauthorizedResponse(w, r, "Invalid token")
			return
		}

		err = app.services.Auth.SetTokenAsUsed(r.Context(), tokenCookie.Value)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) enrichSentryHub(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		transaction := sentry.TransactionFromContext(r.Context())
		transaction.Status = sentry.HTTPtoSpanStatus(rw.statusCode)
	})
}
