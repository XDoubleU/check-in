package main

import (
	"log"
	"os"

	"check-in/api/internal/config"
	"check-in/api/internal/database"
	"check-in/api/internal/services"
)

type application struct {
	config   config.Config
	logger   *log.Logger
	services services.Services
}

//	@title			Check-In API
//	@version		1.0
//	@license.name	GPL-3.0
//	@Accept			json
//	@Produce		json

func main() {
	config := config.New()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := database.Connect(
		config.DB.Dsn,
		config.DB.MaxConns,
		config.DB.MaxIdleTime,
	)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	spandb := database.SpanDB{
		DB: db,
	}

	logger.Printf("connected to database")

	app := &application{
		config:   config,
		logger:   logger,
		services: services.New(spandb),
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
