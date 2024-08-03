package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/xdoubleu/essentia/pkg/database/postgres"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

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

func (repo LocationRepository) GetByIDs(ctx context.Context, ids []string) ([]*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		WHERE locations.id IN $1
	`

	var idsQuery string
	for _, id := range ids {
		idsQuery += fmt.Sprintf("%s,", id)
	}
	// remove last ,
	idsQuery = idsQuery[:len(idsQuery)-1]

	rows, err := repo.db.Query(ctx, query, idsQuery)
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

func (repo LocationRepository) GetByID(ctx context.Context, id string) (*models.Location, error) {
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

func (repo LocationRepository) GetByUserID(ctx context.Context, id string) (*models.Location, error) {
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

// todo: refactor, need tx but don't want code duplication
func (repo LocationRepository) Update(
	ctx context.Context,
	location *models.Location,
	user *models.User,
	updateLocationDto dtos.UpdateLocationDto,
) error {
	locationChanged := false
	userChanged := false

	if updateLocationDto.Name != nil {
		locationChanged = true

		location.Name = *updateLocationDto.Name
	}

	if updateLocationDto.Capacity != nil {
		locationChanged = true

		diff := *updateLocationDto.Capacity - location.Capacity
		location.Available += diff

		if location.Available < 0 {
			location.Available = 0
		}

		location.Capacity = *updateLocationDto.Capacity
	}

	if updateLocationDto.TimeZone != nil {
		locationChanged = true

		location.TimeZone = *updateLocationDto.TimeZone
	}

	if updateLocationDto.Username != nil {
		userChanged = true

		user.Username = *updateLocationDto.Username
	}

	if updateLocationDto.Password != nil {
		userChanged = true

		passwordHash, _ := models.HashPassword(*updateLocationDto.Password)
		user.PasswordHash = passwordHash
	}

	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	if locationChanged {
		err = updateLocation(ctx, tx, location)
		if err != nil {
			return err
		}
	}

	if userChanged {
		err = updateUser(ctx, tx, user)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func updateLocation(ctx context.Context, tx pgx.Tx, location *models.Location) error {
	queryLocation := `
			UPDATE locations
			SET name = $2, capacity = $3, time_zone = $4
			WHERE id = $1
		`

	resultLocation, err := tx.Exec(
		ctx,
		queryLocation,
		location.ID,
		location.Name,
		location.Capacity,
		location.TimeZone,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := resultLocation.RowsAffected()
	if rowsAffected == 0 {
		return errortools.ErrResourceNotFound
	}

	return nil
}

func updateUser(ctx context.Context, tx pgx.Tx, user *models.User) error {
	queryUser := `
			UPDATE users
			SET username = $2, password_hash = $3
			WHERE id = $1 AND role = 'default'
		`

	resultUser, err := tx.Exec(
		ctx,
		queryUser,
		user.ID,
		user.Username,
		user.PasswordHash,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := resultUser.RowsAffected()
	if rowsAffected == 0 {
		return errortools.ErrResourceNotFound
	}

	return nil
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
		return errortools.ErrResourceNotFound
	}

	return nil
}
