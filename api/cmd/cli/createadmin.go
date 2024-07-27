//nolint:forbidigo //returns output of cli
package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/xdoubleu/essentia/pkg/database/postgres"

	"check-in/api/internal/config"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
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
