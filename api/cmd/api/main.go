package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
	_ "time/tzdata"

	"github.com/xdoubleu/essentia/pkg/database/postgres"
	"github.com/xdoubleu/essentia/pkg/httptools"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/sentrytools"

	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
	"check-in/api/internal/services"
)

type Application struct {
	logger   *slog.Logger
	config   config.Config
	services services.Services
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

//	@title			Check-In API
//	@version		1.0
//	@license.name	GPL-3.0
//	@Accept			json
//	@Produce		json

func main() {
	cfg := config.New()

	logger := slog.New(sentrytools.NewSentryLogHandler())
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
