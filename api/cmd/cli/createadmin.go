//nolint:forbidigo //returns output of cli
package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"check-in/api/internal/config"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
	"check-in/api/internal/services"
)

func createAdmin(cfg config.Config, username string, password string) {
	if username == "" || password == "" {
		fmt.Println("please provide a username and password")
		return
	}

	db, err := postgres.Connect(
		slog.Default(),
		cfg.DBDsn,
		25, //nolint:mnd //no magic number
		"15m",
		10,             //nolint:mnd //no magic number
		10*time.Second, //nolint:mnd //no magic number
		30*time.Second, //nolint:mnd //no magic number
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	services := services.New(
		context.Background(),
		slog.Default(),
		cfg,
		repositories.New(db),
	)

	_, err = services.Users.Create(
		context.Background(),
		//nolint:exhaustruct //other fields are optional
		&dtos.CreateUserDto{
			Username: username,
			Password: password,
		},
		models.AdminRole,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("admin added")
}
