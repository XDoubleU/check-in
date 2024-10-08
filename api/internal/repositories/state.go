package repositories

import (
	"context"
	"strconv"
	"time"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"check-in/api/internal/models"
)

type StateRepository struct {
	db postgres.DB
}

func (repo StateRepository) Get(
	ctx context.Context,
) (*models.State, error) {
	query := `
		SELECT key, value
		FROM states
	`

	//nolint:exhaustruct //fields are initialized later
	state := models.State{}

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	for rows.Next() {
		var key models.StateKey
		var value string

		err = rows.Scan(
			&key,
			&value,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		if key == models.IsMaintenanceKey {
			state.IsMaintenance, err = strconv.ParseBool(value)
		}

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &state, nil
}

func (repo StateRepository) IsDatabaseActive(ctx context.Context) bool {
	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return repo.db.Ping(ctx2) == nil
}

func (repo StateRepository) UpdateKey(
	ctx context.Context,
	key models.StateKey,
	value string,
) error {
	query := `
		UPDATE states
		SET value = $2
		WHERE key = $1
	`

	result, err := repo.db.Exec(ctx, query, key, value)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}
