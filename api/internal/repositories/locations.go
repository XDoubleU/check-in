package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xdoubleu/essentia/pkg/database"
	"github.com/xdoubleu/essentia/pkg/database/postgres"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type LocationRepository struct {
	db postgres.DB
}

func (repo LocationRepository) GetTotalCount(ctx context.Context) (*int64, error) {
	query := `
		SELECT COUNT(*)
		FROM locations
	`

	var total *int64

	err := repo.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return total, nil
}

func (repo LocationRepository) GetAll(ctx context.Context) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		ORDER BY name ASC
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	locations := []*models.Location{}
	for rows.Next() {
		var location models.Location

		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.Capacity,
			&location.TimeZone,
			&location.UserID,
		)
		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return locations, nil
}

func (repo LocationRepository) GetAllPaginated(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	locations := []*models.Location{}
	for rows.Next() {
		var location models.Location

		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.Capacity,
			&location.TimeZone,
			&location.UserID,
		)
		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return locations, nil
}

func (repo LocationRepository) GetByIDs(
	ctx context.Context,
	ids []string,
) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		WHERE locations.id = ANY($1::uuid[])
	`

	if len(ids) == 0 {
		return make([]*models.Location, 0), nil
	}

	if len(ids) == 1 {
		location, err := repo.GetByID(ctx, ids[0])
		if err != nil {
			return nil, err
		}

		return []*models.Location{location}, nil
	}

	//nolint:exhaustruct //other fields are optional
	pgArray := pgtype.Array[pgtype.Text]{}
	for _, id := range ids {
		pgArray.Elements = append(pgArray.Elements, pgtype.Text{
			String: id,
			Valid:  true,
		})
	}

	rows, err := repo.db.Query(ctx, query, pgArray)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	locations := []*models.Location{}
	for rows.Next() {
		var location models.Location

		err = rows.Scan(
			&location.ID,
			&location.Name,
			&location.Capacity,
			&location.TimeZone,
			&location.UserID,
		)
		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return locations, nil
}

func (repo LocationRepository) GetByID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		WHERE locations.id = $1
	`

	//nolint:exhaustruct //other fields are optional
	location := models.Location{}
	err := repo.db.QueryRow(
		ctx,
		query,
		id).Scan(
		&location.ID,
		&location.Name,
		&location.Capacity,
		&location.TimeZone,
		&location.UserID,
	)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &location, nil
}

func (repo LocationRepository) GetByUserID(
	ctx context.Context,
	id string,
) (*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		WHERE user_id = $1
	`

	//nolint:exhaustruct //other fields are optional
	location := models.Location{}
	err := repo.db.QueryRow(
		ctx,
		query,
		id).Scan(
		&location.ID,
		&location.Name,
		&location.Capacity,
		&location.TimeZone,
		&location.UserID,
	)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &location, nil
}

func (repo LocationRepository) Create(
	ctx context.Context,
	name string,
	capacity int64,
	timeZone string,
	userID string,
) (*models.Location, error) {
	query := `
		INSERT INTO locations (name, capacity, time_zone, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	//nolint:exhaustruct //other fields are optional
	location := models.Location{
		Name:      name,
		Capacity:  capacity,
		Available: capacity,
		TimeZone:  timeZone,
		UserID:    userID,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		name,
		capacity,
		timeZone,
		userID,
	).Scan(&location.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &location, nil
}

func (repo LocationRepository) Update(
	ctx context.Context,
	location models.Location,
	updateLocationDto *dtos.UpdateLocationDto,
) (*models.Location, error) {
	query := `
		UPDATE locations
		SET name = $2, capacity = $3, time_zone = $4
		WHERE id = $1
	`

	if updateLocationDto.Name != nil {
		location.Name = *updateLocationDto.Name
	}

	if updateLocationDto.Capacity != nil {
		diff := *updateLocationDto.Capacity - location.Capacity
		location.Available += diff

		if location.Available < 0 {
			location.Available = 0
		}

		location.Capacity = *updateLocationDto.Capacity
	}

	if updateLocationDto.TimeZone != nil {
		location.TimeZone = *updateLocationDto.TimeZone
	}

	resultLocation, err := repo.db.Exec(
		ctx,
		query,
		location.ID,
		location.Name,
		location.Capacity,
		location.TimeZone,
	)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := resultLocation.RowsAffected()
	if rowsAffected == 0 {
		return nil, database.ErrResourceNotFound
	}

	return &location, nil
}

func (repo LocationRepository) Delete(
	ctx context.Context,
	location *models.Location,
) error {
	query := `
		DELETE FROM locations
		WHERE id = $1
	`

	result, err := repo.db.Exec(ctx, query, location.ID)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}
