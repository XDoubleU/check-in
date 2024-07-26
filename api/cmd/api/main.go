package main

import (
	"log/slog"
	_ "time/tzdata"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/logger"

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
	slog.Info(cfg.String())

	spandb := postgres.SpanDB{
		DB: db,
	}

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

	db, err := postgres.Connect(
		cfg.DB.Dsn,
		cfg.DB.MaxConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := NewApp(slog.Default(), cfg, db)

	err = httptools.Serve(app.config.Port, app.routes(), app.config.Env)
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
}
