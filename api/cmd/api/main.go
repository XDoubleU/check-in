package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/logging"
	sentrytools "github.com/XDoubleU/essentia/pkg/sentry"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
	"check-in/api/internal/services"
	"check-in/api/internal/shared"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Application struct {
	logger     *slog.Logger
	ctx        context.Context
	ctxCancel  context.CancelFunc
	db         postgres.DB
	config     config.Config
	getTimeNow shared.NowTimeProvider
	services   services.Services
}

//	@title			Check-In API
//	@version		1.0
//	@license.name	GPL-3.0
//	@Accept			json
//	@Produce		json

func main() {
	cfg := config.New()

	logger := slog.New(sentrytools.NewLogHandler(cfg.Env,
		slog.NewTextHandler(os.Stdout, nil)))
	db, err := postgres.Connect(
		logger,
		cfg.DBDsn,
		25, //nolint:mnd //no magic number
		"15m",
		30,             //nolint:mnd //no magic number
		30*time.Second, //nolint:mnd //no magic number
		5*time.Minute,  //nolint:mnd //no magic number
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ApplyMigrations(logger, db)

	app := NewApp(logger, cfg, db, time.Now)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,  //nolint:mnd //no magic number
		WriteTimeout: 10 * time.Second, //nolint:mnd //no magic number
	}
	err = httptools.Serve(logger, srv, app.config.Env)
	if err != nil {
		logger.Error("failed to serve server", logging.ErrAttr(err))
	}
}

func NewApp(
	logger *slog.Logger,
	cfg config.Config,
	db postgres.DB,
	nowTimeProvider shared.NowTimeProvider,
) *Application {
	logger.Info(cfg.String())

	//nolint:exhaustruct //other fields are optional
	app := &Application{
		logger:     logger,
		config:     cfg,
		getTimeNow: nowTimeProvider,
	}

	app.setContext()
	app.setDB(db)

	return app
}

func (app *Application) getNowTimeProvider() shared.NowTimeProvider {
	return func() time.Time { return app.getTimeNow() }
}

func (app *Application) setDB(db postgres.DB) {
	// make sure previous app is cancelled internally
	app.ctxCancel()

	app.setContext()

	spandb := postgres.NewSpanDB(db)

	app.db = spandb
	app.services = services.New(
		app.ctx,
		app.logger,
		app.config,
		repositories.New(app.db, app.getNowTimeProvider()),
		app.getNowTimeProvider(),
	)
}

func (app *Application) setContext() {
	ctx, cancel := context.WithCancel(context.Background())
	app.ctx = ctx
	app.ctxCancel = cancel
}

func ApplyMigrations(logger *slog.Logger, db *pgxpool.Pool) {
	migrationsDB := stdlib.OpenDBFromPool(db)

	goose.SetLogger(slog.NewLogLogger(logger.Handler(), slog.LevelInfo))

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		panic(err)
	}

	if err := goose.Up(migrationsDB, "migrations"); err != nil {
		panic(err)
	}
}
