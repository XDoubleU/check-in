package services

import (
	"context"

	"check-in/api/internal/database"
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type SchoolService struct {
	db database.DB
}

func (service SchoolService) GetTotalCount(ctx context.Context) (*int64, error) {
	query := `
		SELECT COUNT(*)
		FROM schools
	`

	var total *int64

	err := service.db.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return nil, handleError(err)
	}

	return total, nil
}

func (service SchoolService) GetAll(ctx context.Context) ([]*models.School, error) {
	query := `
		SELECT id, name
		FROM schools
		ORDER BY name ASC
	`

	rows, err := service.db.Query(ctx, query)
	if err != nil {
		return nil, handleError(err)
	}

	schools := []*models.School{}

	for rows.Next() {
		var school models.School

		err = rows.Scan(
			&school.ID,
			&school.Name,
		)

		if err != nil {
			return nil, handleError(err)
		}

		schools = append(schools, &school)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return schools, nil
}

func (service SchoolService) GetAllSortedByLocation(
	ctx context.Context,
	locationID string,
) ([]*models.School, error) {
	query := `
		SELECT id, name
		FROM schools
		ORDER BY
			CASE
				WHEN read_only = true THEN -1
				ELSE (
					SELECT COUNT(*)
					FROM check_ins
					WHERE check_ins.location_id = $1
					AND check_ins.school_id = schools.id
				)
			END
		DESC, name ASC
	`

	rows, err := service.db.Query(ctx, query, locationID)
	if err != nil {
		return nil, handleError(err)
	}

	schools := []*models.School{}

	for rows.Next() {
		var school models.School

		err = rows.Scan(
			&school.ID,
			&school.Name,
		)

		if err != nil {
			return nil, handleError(err)
		}

		schools = append(schools, &school)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return schools, nil
}

func (service SchoolService) GetAllPaginated(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]*models.School, error) {
	query := `
		SELECT id, name, read_only
		FROM schools
		ORDER BY name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := service.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, handleError(err)
	}

	schools := []*models.School{}

	for rows.Next() {
		var school models.School

		err = rows.Scan(
			&school.ID,
			&school.Name,
			&school.ReadOnly,
		)

		if err != nil {
			return nil, handleError(err)
		}

		schools = append(schools, &school)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return schools, nil
}

func (service SchoolService) GetByID(
	ctx context.Context,
	id int64,
) (*models.School, error) {
	query := `
		SELECT name, read_only
		FROM schools
		WHERE id = $1
	`

	school := models.School{
		ID: id,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		id).Scan(&school.Name, &school.ReadOnly)

	if err != nil {
		return nil, handleError(err)
	}

	return &school, nil
}

func (service SchoolService) GetByIDWithoutReadOnly(
	ctx context.Context,
	id int64,
) (*models.School, error) {
	query := `
		SELECT name
		FROM schools
		WHERE id = $1 AND read_only = false
	`

	school := models.School{
		ID:       id,
		ReadOnly: false,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		id).Scan(&school.Name)

	if err != nil {
		return nil, handleError(err)
	}

	return &school, nil
}

func (service SchoolService) Create(
	ctx context.Context,
	name string,
) (*models.School, error) {
	query := `
		INSERT INTO schools (name)
		VALUES ($1)
		RETURNING id
	`

	school := models.School{
		Name: name,
	}

	err := service.db.QueryRow(ctx, query, name).Scan(&school.ID)

	if err != nil {
		return nil, handleError(err)
	}

	return &school, nil
}

func (service SchoolService) Update(
	ctx context.Context,
	school *models.School,
	schoolDto dtos.SchoolDto,
) error {
	school.Name = schoolDto.Name

	query := `
		UPDATE schools
		SET name = $2
		WHERE id = $1 AND read_only = false
	`

	result, err := service.db.Exec(ctx, query, school.ID, school.Name)
	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (service SchoolService) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM schools
		WHERE id = $1 AND read_only = false
	`

	result, err := service.db.Exec(ctx, query, id)
	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
