package repositories

import (
	"context"
	"time"

	"github.com/xdoubleu/essentia/pkg/database/postgres"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/models"
)

type CheckInRepository struct {
	db postgres.DB
}

func (repo CheckInRepository) GetAllInRange(
	ctx context.Context,
	locationIDs []string,
	startDate time.Time,
	endDate time.Time,
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

	rows, err := repo.db.Query(
		ctx,
		query,
		locationIDs,
		startDate,
		endDate,
	)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
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
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		checkIns = append(checkIns, &checkIn)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return checkIns, nil
}

func (repo CheckInRepository) GetByID(
	ctx context.Context,
	location *models.Location,
	id int64,
) (*models.CheckIn, error) {
	query := `
		SELECT school_id, capacity, created_at AT TIME ZONE $3
		FROM check_ins
		WHERE id = $1 AND location_id = $2
	`

	//nolint:exhaustruct //other fields are optional
	checkIn := models.CheckIn{
		ID:         id,
		LocationID: location.ID,
	}

	err := repo.db.QueryRow(
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
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &checkIn, nil
}

func (repo CheckInRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM check_ins
		WHERE id = $1
	`

	result, err := repo.db.Exec(ctx, query, id)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errortools.ErrResourceNotFound
	}

	return nil
}
