//nolint:forbidigo //returns output of cli
package main

import (
	"context"
	"fmt"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"check-in/api/internal/config"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

func createAdmin(cfg config.Config, username string, password string) {
	if username == "" || password == "" {
		fmt.Println("please provide a username and password")
		return
	}

	db, err := postgres.Connect(cfg.DB.Dsn, cfg.DB.MaxConns, cfg.DB.MaxIdleTime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	srvs := repositories.New(db)

	_, err = srvs.Users.Create(
		context.Background(),
		username,
		password,
		models.AdminRole,
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("admin added")
}
