package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dsn string, maxConns int, maxIdleTime string) (*pgxpool.Pool, error) {
	connString, err := editConnString(dsn, maxConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second, //nolint:gomnd //no magic number
		)
		defer cancel()

		err = db.Ping(ctx)

		if err == nil || i == 2 {
			break
		}

		retryTime := 15 * time.Second //nolint:gomnd //no magic number
		fmt.Printf(                   //nolint:forbidigo //allowed printf
			"can't connect to database, retrying in %s",
			retryTime,
		)
		time.Sleep(retryTime)
	}

	if err != nil {
		return nil, errors.New("can't connect to database")
	}

	return db, nil
}

func editConnString(dsn string, maxConns int, maxIdleTime string) (string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	queryValues := parsedURL.Query()

	queryValues.Add("pool_max_conns", strconv.Itoa(maxConns))
	queryValues.Add("pool_max_conn_idle_time", maxIdleTime)

	parsedURL.RawQuery = queryValues.Encode()

	return parsedURL.String(), nil
}
