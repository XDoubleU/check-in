//nolint:forbidigo //returns output of cli
package main

import (
	"flag"
	"fmt"

	"check-in/api/internal/config"
)

func main() {
	cfg := config.New()
	var username string
	var password string

	flag.StringVar(
		&cfg.DBDsn,
		"db",
		"postgres://postgres@localhost/postgres",
		"DB DSN",
	)
	flag.StringVar(&username, "u", "", "username of admin user")
	flag.StringVar(&password, "p", "", "password of admin user")

	flag.Parse()

	command := flag.Arg(0)
	switch command {
	case "createadmin":
		createAdmin(cfg, username, password)
	default:
		fmt.Println("invalid command")
		return
	}
}
