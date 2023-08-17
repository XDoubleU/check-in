package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func (app *application) getSentryHandler() *sentryhttp.Handler {
	if len(app.config.SentryDsn) == 0 {
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              app.config.SentryDsn,
		Environment:      app.config.Env,
		Release:          app.config.Release,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	})

	if err != nil {
		app.logger.Printf("sentry initialization failed: %v\n", err)
		return nil
	}

	return sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
}

func sentryGoRoutineErrorHandler(name string, f func(ctx context.Context) error) {
	name = fmt.Sprintf("GO ROUTINE %s", name)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second, //nolint:gomnd // no magic number
	)
	defer cancel()

	hub := sentry.CurrentHub().Clone()
	ctx = sentry.SetHubOnContext(ctx, hub)

	options := []sentry.SpanOption{
		sentry.WithOpName("go.routine"),
	}

	transaction := sentry.StartTransaction(ctx, name, options...)
	transaction.Status = sentry.HTTPtoSpanStatus(http.StatusOK)
	defer transaction.Finish()

	err := f(transaction.Context())

	if err != nil {
		transaction.Status = sentry.HTTPtoSpanStatus(http.StatusInternalServerError)

		hub.WithScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelError)
			hub.CaptureException(err)
		})
	}
}
