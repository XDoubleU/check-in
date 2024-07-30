package main

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	_ "time/tzdata"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	httptools "github.com/xdoubleu/essentia/pkg/communication/http"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	"github.com/xdoubleu/essentia/pkg/logging"
	sentrytools "github.com/xdoubleu/essentia/pkg/sentry"

	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
	"check-in/api/internal/services"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Application struct {
	logger   *slog.Logger
	config   config.Config
	services services.Services
}

//	@title			Check-In API
//	@version		1.0
//	@license.name	GPL-3.0
//	@Accept			json
//	@Produce		json

func main() {
	cfg := config.New()

	logger := slog.New(sentrytools.NewLogHandler())
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

	app := NewApp(logger, cfg, db)

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

func NewApp(logger *slog.Logger, cfg config.Config, db postgres.DB) *Application {
	logger.Info(cfg.String())

	spandb := postgres.NewSpanDB(db)

	repos := repositories.New(spandb)

	return &Application{
		logger:   logger,
		config:   cfg,
		services: services.New(cfg, repos),
	}
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
