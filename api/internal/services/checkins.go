package services

import (
	"context"
	"time"

	"check-in/api/internal/helpers"
	"check-in/api/internal/models"

	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/XDoubleU/essentia/pkg/http_tools"
)

type CheckInService struct {
	db postgres.DB
}

func (service CheckInService) GetAllOfDay(
	ctx context.Context,
	locationID string,
	date *time.Time,
) ([]*models.CheckIn, error) {
	return service.GetAllInRange(
		ctx,
		[]string{locationID},
		helpers.StartOfDay(date),
		helpers.EndOfDay(date),
	)
}

func (service CheckInService) GetAllInRange(
	ctx context.Context,
	locationIDs []string,
	startDate *time.Time,
	endDate *time.Time,
) ([]*models.CheckIn, error) {
	query := `
		SELECT check_ins.id, check_ins.location_id, check_ins.school_id,
		 check_ins.capacity, (check_ins.created_at AT TIME ZONE locations.time_zone)
		FROM check_ins
		INNER JOIN locations
		ON locations.id = check_ins.location_id
		WHERE check_ins.location_id = any($1)
		AND (check_ins.created_at AT TIME ZONE locations.time_zone) >= $2
		AND (check_ins.created_at AT TIME ZONE locations.time_zone) <= $3
		ORDER BY check_ins.created_at 
	`

	rows, err := service.db.Query(
		ctx,
		query,
		locationIDs,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, handleError(err)
	}

	checkIns := []*models.CheckIn{}

	for rows.Next() {
		var checkIn models.CheckIn

		err = rows.Scan(
			&checkIn.ID,
			&checkIn.LocationID,
			&checkIn.SchoolID,
			&checkIn.Capacity,
			&checkIn.CreatedAt,
		)

		if err != nil {
			return nil, handleError(err)
		}

		checkIns = append(checkIns, &checkIn)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError(err)
	}

	return checkIns, nil
}

func (service CheckInService) GetByID(
	ctx context.Context,
	location *models.Location,
	id int64,
) (*models.CheckIn, error) {
	query := `
		SELECT school_id, capacity, created_at AT TIME ZONE $3
		FROM check_ins
		WHERE id = $1 AND location_id = $2
	`

	checkIn := models.CheckIn{
		ID:         id,
		LocationID: location.ID,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		id,
		location.ID,
		location.TimeZone,
	).Scan(
		&checkIn.SchoolID,
		&checkIn.Capacity,
		&checkIn.CreatedAt,
	)

	if err != nil {
		return nil, handleError(err)
	}

	return &checkIn, nil
}

func (service CheckInService) Create(
	ctx context.Context,
	location *models.Location,
	school *models.School,
) (*models.CheckIn, error) {
	query := `
		INSERT INTO check_ins (location_id, school_id, capacity)
		VALUES ($1, $2, $3)
		RETURNING id, (created_at AT TIME ZONE $4)
	`

	checkIn := models.CheckIn{
		LocationID: location.ID,
		SchoolID:   school.ID,
		Capacity:   location.Capacity,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		location.ID,
		school.ID,
		location.Capacity,
		location.TimeZone,
	).Scan(&checkIn.ID, &checkIn.CreatedAt)

	if err != nil {
		return nil, handleError(err)
	}

	location.Available--

	return &checkIn, nil
}

func (service CheckInService) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM check_ins
		WHERE id = $1
	`

	result, err := service.db.Exec(ctx, query, id)
	if err != nil {
		return handleError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return http_tools.ErrRecordNotFound
	}

	return nil
}
