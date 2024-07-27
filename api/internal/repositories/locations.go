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

// todo: refactor
func (repo LocationRepository) GetBy(
	ctx context.Context,
	whereQuery string,
	value string,
) (*models.Location, error) {
	query := `
		SELECT id, name, capacity, time_zone, user_id
		FROM locations
		%s
	`

	query = fmt.Sprintf(query, whereQuery)

	//nolint:exhaustruct //other fields are optional
	location := models.Location{}

	err := repo.db.QueryRow(
		ctx,
		query,
		value).Scan(
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

// todo: refactor, need tx but don't want
// code duplication across repositories
func (repo LocationRepository) Create(
	ctx context.Context,
	name string,
	capacity int64,
	timeZone string,
	username string,
	password string,
) (*models.Location, error) {
	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	var user *models.User
	var location *models.Location

	user, err = createUser(ctx, tx, username, password)
	if err != nil {
		return nil, err
	}

	location, err = createLocation(ctx, tx, name, capacity, timeZone, user.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func createUser(
	ctx context.Context,
	tx pgx.Tx,
	username string,
	password string,
) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, 'default')
		RETURNING id
	`

	//nolint:exhaustruct //other fields are optional
	user := models.User{
		Username: username,
		Role:     models.DefaultRole,
	}

	passwordHash, err := models.HashPassword(password)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(
		ctx,
		query,
		username,
		passwordHash,
	).Scan(&user.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &user, nil
}

func createLocation(
	ctx context.Context,
	tx pgx.Tx,
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

	err := tx.QueryRow(
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

// todo: refactor, need tx but don't want code duplication
func (repo LocationRepository) Delete(
	ctx context.Context,
	location *models.Location,
) error {
	tx, err := repo.db.Begin(ctx)
	defer tx.Rollback(ctx) //nolint:errcheck //deferred
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	err = deleteLocation(ctx, tx, location.ID)
	if err != nil {
		return err
	}

	err = deleteUser(ctx, tx, location.UserID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func deleteLocation(ctx context.Context, tx pgx.Tx, id string) error {
	query := `
		DELETE FROM locations
		WHERE id = $1
	`

	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errortools.ErrResourceNotFound
	}

	return nil
}

func deleteUser(ctx context.Context, tx pgx.Tx, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1 AND role = 'default'
	`

	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errortools.ErrResourceNotFound
	}

	return nil
}
