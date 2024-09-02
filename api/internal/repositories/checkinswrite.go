package repositories

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"check-in/api/internal/models"
)

type CheckInWriteRepository struct {
	db postgres.DB
}

func (repo CheckInWriteRepository) Create(
	ctx context.Context,
	location *models.Location,
	school *models.School,
) (*models.CheckIn, error) {
	query := `
		INSERT INTO check_ins (location_id, school_id, capacity)
		VALUES ($1, $2, $3)
		RETURNING id, (created_at AT TIME ZONE $4)
	`

	//nolint:exhaustruct //other fields are optional
	checkIn := models.CheckIn{
		LocationID: location.ID,
		SchoolID:   school.ID,
		Capacity:   location.Capacity,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		location.ID,
		school.ID,
		location.Capacity,
		location.TimeZone,
	).Scan(&checkIn.ID, &checkIn.CreatedAt)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &checkIn, nil
}
