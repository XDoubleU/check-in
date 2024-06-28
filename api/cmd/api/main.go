package main

import (
	_ "time/tzdata"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/httptools"
	"github.com/XDoubleU/essentia/pkg/logger"

	"check-in/api/internal/config"
	"check-in/api/internal/repositories"
)

type application struct {
	config       config.Config
	repositories repositories.Repositories
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

	spandb := postgres.SpanDB{
		DB: db,
	}

	logger.GetLogger().Printf("connected to database")

	app := &application{
		config:       cfg,
		repositories: repositories.New(spandb),
	}

	app.config.Print()

	err = httptools.Serve(app.config.Port, app.routes(), app.config.Env)
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
}
