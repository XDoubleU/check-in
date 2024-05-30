package services

import (
	"context"

	"check-in/api/internal/database"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"

	"github.com/XDoubleU/essentia/pkg/http_tools"
)

type UserService struct {
	db        database.DB
	locations LocationService
}

func (service UserService) GetTotalCount(ctx context.Context) (*int64, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE role = 'manager'
	`

	var total *int64

	err := service.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return nil, handleError(err)
	}

	return total, nil
}

func (service UserService) GetAll(
	ctx context.Context,
) ([]*models.User, error) {
	query := `
		SELECT id, username
		FROM users
		WHERE role = 'manager'
	`

	rows, err := service.db.Query(ctx, query)
	if err != nil {
		return nil, handleError(err)
	}

	users := []*models.User{}

	for rows.Next() {
		user := models.User{
			Role: models.ManagerRole,
		}

		err = rows.Scan(
			&user.ID,
			&user.Username,
		)

		if err != nil {
			return nil, handleError(err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return users, nil
}

func (service UserService) GetAllPaginated(
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

	rows, err := service.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, handleError(err)
	}

	users := []*models.User{}

	for rows.Next() {
		user := models.User{
			Role: models.ManagerRole,
		}

		err = rows.Scan(
			&user.ID,
			&user.Username,
		)

		if err != nil {
			return nil, handleError(err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return users, nil
}

func (service UserService) GetByID(
	ctx context.Context,
	id string,
	role models.Role,
) (*models.User, error) {
	query := `
		SELECT users.username, users.password_hash
		FROM users
		WHERE users.id = $1 AND users.role = $2
	`

	user := models.User{
		ID:   id,
		Role: role,
	}
	err := service.db.QueryRow(
		ctx,
		query,
		id,
		role,
	).Scan(&user.Username, &user.PasswordHash)

	if err != nil {
		return nil, handleError(err)
	}

	if user.Role == models.DefaultRole {
		var location *models.Location
		location, err = service.locations.GetByUserID(ctx, user.ID)
		if err != nil {
			return nil, handleError(err)
		}

		user.Location = location
	}

	return &user, nil
}

func (service UserService) GetByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	query := `
		SELECT id, password_hash, role
		FROM users
		WHERE username = $1
	`

	user := models.User{
		Username: username,
	}
	err := service.db.QueryRow(
		ctx,
		query,
		username,
	).Scan(&user.ID, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, handleError(err)
	}

	return &user, nil
}

func (service UserService) Create(
	ctx context.Context,
	username string,
	password string,
	role models.Role,
) (*models.User, error) {
	query := `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	user := models.User{
		Username: username,
		Role:     role,
	}

	passwordHash, err := models.HashPassword(password)
	if err != nil {
		return nil, err
	}

	err = service.db.QueryRow(
		ctx,
		query,
		username,
		passwordHash,
		role,
	).Scan(&user.ID)

	if err != nil {
		return nil, handleError(err)
	}

	return &user, nil
}

func (service UserService) Update(
	ctx context.Context,
	user *models.User,
	updateUserDto dtos.UpdateUserDto,
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

	result, err := service.db.Exec(
		ctx,
		query,
		user.ID,
		role,
		user.Username,
		user.PasswordHash,
	)

	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return http_tools.ErrRecordNotFound
	}

	return nil
}

func (service UserService) Delete(
	ctx context.Context,
	id string,
	role models.Role,
) error {
	query := `
		DELETE FROM users
		WHERE id = $1 AND role = $2
	`

	result, err := service.db.Exec(ctx, query, id, role)
	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return http_tools.ErrRecordNotFound
	}

	return nil
}
