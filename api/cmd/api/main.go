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

type application struct {
	logger   *slog.Logger
	config   config.Config
	services services.Services
}

func NewApp(logger *slog.Logger, cfg config.Config, db postgres.DB) *application {
	logger.Info(cfg.String())

	spandb := postgres.NewSpanDB(db)

	logger.Info("connected to database")

	repos := repositories.New(spandb)

	return &application{
		logger:   logger,
		config:   config.New(),
		services: services.New(repos),
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
		25,
		"15m",
		30,
		30*time.Second,
		5*time.Minute,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := NewApp(logger, cfg, db)

	routes, err := app.routes()
	if err != nil {
		logger.Error("failed to setup routes", logging.ErrAttr(err))
		return
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      *routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,  //nolint:gomnd //no magic number
		WriteTimeout: 10 * time.Second, //nolint:gomnd //no magic number
	}
	err = httptools.Serve(logger, srv, app.config.Env)
	if err != nil {
		logger.Error("failed to serve server", logging.ErrAttr(err))
	}
}
