//nolint:forbidigo //returns output of cli
package main

import (
	"context"
	"fmt"

	"check-in/api/internal/config"
	"check-in/api/internal/database"
	"check-in/api/internal/models"
	"check-in/api/internal/services"
)

func createAdmin(cfg config.Config, username string, password string) {
	if username == "" || password == "" {
		fmt.Println("please provide a username and password")
		return
	}

	db, err := database.Connect(cfg.DB.Dsn, cfg.DB.MaxConns, cfg.DB.MaxIdleTime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	srvs := services.New(db)

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
