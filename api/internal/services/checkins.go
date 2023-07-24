package services

import (
	"context"
	"time"

	"check-in/api/internal/database"
	"check-in/api/internal/models"
)

type CheckInService struct {
	db database.DB
}

func (service CheckInService) GetAllInRange(
	ctx context.Context,
	locationID string,
	startDate *time.Time,
	endDate *time.Time,
) ([]*models.CheckIn, error) {
	query := `
		SELECT school_id, capacity, created_at
		FROM check_ins
		WHERE location_id = $1
		AND created_at >= $2
		AND created_at <= $3
	`

	rows, err := service.db.Query(ctx, query, locationID, startDate, endDate)
	if err != nil {
		return nil, handleError(err)
	}

	checkIns := []*models.CheckIn{}

	for rows.Next() {
		var checkIn models.CheckIn

		err = rows.Scan(
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

func (service CheckInService) Create(
	ctx context.Context,
	locationID string,
	schoolID int64,
	capacity int64,
) (*models.CheckIn, error) {
	query := `
		INSERT INTO check_ins (location_id, school_id, capacity)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	checkIn := models.CheckIn{
		LocationID: locationID,
		SchoolID:   schoolID,
		Capacity:   capacity,
	}

	err := service.db.QueryRow(
		ctx,
		query,
		locationID,
		schoolID,
		capacity,
	).Scan(&checkIn.ID, &checkIn.CreatedAt)

	if err != nil {
		return nil, handleError(err)
	}

	return &checkIn, nil
}
