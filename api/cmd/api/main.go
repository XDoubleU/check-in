package main

import (
	_ "time/tzdata"

	"check-in/api/internal/config"
	"check-in/api/internal/database"
	"check-in/api/internal/services"

	"github.com/XDoubleU/essentia/pkg/http_tools"
)

type application struct {
	config     config.Config
	services   services.Services
	hideErrors bool
}

//	@title			Check-In API
//	@version		1.0
//	@license.name	GPL-3.0
//	@Accept			json
//	@Produce		json

func main() {
	cfg := config.New()

	db, err := database.Connect(
		cfg.DB.Dsn,
		cfg.DB.MaxConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	spandb := database.SpanDB{
		DB: db,
	}

	http_tools.GetLogger().Printf("connected to database")

	app := &application{
		config:     cfg,
		services:   services.New(spandb),
		hideErrors: cfg.Env == config.ProdEnv,
	}

	app.config.Print()

	err = http_tools.Serve(app.config.Port, app.routes(), app.config.Env)
	if err != nil {
		http_tools.GetLogger().Fatal(err)
	}
}
