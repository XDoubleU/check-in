package repositories

import (
	"context"

	"github.com/xdoubleu/essentia/pkg/database/postgres"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type UserRepository struct {
	db postgres.DB
}

func (repo UserRepository) GetTotalCount(ctx context.Context) (*int64, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE role = 'manager'
	`

	var total *int64

	err := repo.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return total, nil
}

func (repo UserRepository) GetAll(
	ctx context.Context,
) ([]*models.User, error) {
	query := `
		SELECT id, username
		FROM users
		WHERE role = 'manager'
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	users := []*models.User{}

	for rows.Next() {
		//nolint:exhaustruct //other fields are optional
		user := models.User{
			Role: models.ManagerRole,
		}

		err = rows.Scan(
			&user.ID,
			&user.Username,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return users, nil
}

func (repo UserRepository) GetAllPaginated(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]*models.User, error) {
	query := `
		SELECT id, username
		FROM users
		WHERE role = 'manager'
		ORDER BY username ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	users := []*models.User{}

	for rows.Next() {
		//nolint:exhaustruct //other fields are optional
		user := models.User{
			Role: models.ManagerRole,
		}

		err = rows.Scan(
			&user.ID,
			&user.Username,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return users, nil
}

func (repo UserRepository) GetByID(
	ctx context.Context,
	id string,
	role models.Role,
) (*models.User, error) {
	query := `
		SELECT users.username, users.password_hash
		FROM users
		WHERE users.id = $1 AND users.role = $2
	`

	//nolint:exhaustruct //other fields are optional
	user := models.User{
		ID:   id,
		Role: role,
	}
	err := repo.db.QueryRow(
		ctx,
		query,
		id,
		role,
	).Scan(&user.Username, &user.PasswordHash)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &user, nil
}

func (repo UserRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	query := `
		SELECT id, password_hash, role
		FROM users
		WHERE username = $1
	`

	//nolint:exhaustruct //other fields are optional
	user := models.User{
		Username: username,
	}
	err := repo.db.QueryRow(
		ctx,
		query,
		username,
	).Scan(&user.ID, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &user, nil
}

func (repo UserRepository) Create(
	ctx context.Context,
	username string,
	passwordHash []byte,
	role models.Role,
) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	//nolint:exhaustruct //other fields are optional
	user := models.User{
		Username: username,
		Role:     role,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		username,
		passwordHash,
		role,
	).Scan(&user.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &user, nil
}

func (repo UserRepository) Update(
	ctx context.Context,
	user *models.User,
	updateUserDto *dtos.UpdateUserDto,
	role models.Role,
) error {
	if updateUserDto.Username != nil {
		user.Username = *updateUserDto.Username
	}

	if updateUserDto.Password != nil {
		passwordHash, _ := models.HashPassword(*updateUserDto.Password)
		user.PasswordHash = passwordHash
	}

	query := `
		UPDATE users
		SET username = $3, password_hash = $4
		WHERE id = $1 AND role = $2
	`

	result, err := repo.db.Exec(
		ctx,
		query,
		user.ID,
		role,
		user.Username,
		user.PasswordHash,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errortools.ErrResourceNotFound
	}

	return nil
}

func (repo UserRepository) Delete(
	ctx context.Context,
	id string,
	role models.Role,
) error {
	query := `
		DELETE FROM users
		WHERE id = $1 AND role = $2
	`

	result, err := repo.db.Exec(ctx, query, id, role)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errortools.ErrResourceNotFound
	}

	return nil
}
