package services

import (
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
	"context"
	"time"

	"github.com/XDoubleU/essentia/pkg/tools"
)

type CheckInService struct {
	checkins repositories.CheckInRepository
}

func (service CheckInService) GetAllOfDay(
	ctx context.Context,
	locationID string,
	date time.Time,
) ([]*models.CheckIn, error) {
	return service.checkins.GetAllInRange(
		ctx,
		[]string{locationID},
		tools.StartOfDay(date),
		tools.EndOfDay(date),
	)
}

func (service CheckInService) GetAllInRange(
	ctx context.Context,
	locationIDs []string,
	startDate time.Time,
	endDate time.Time,
) ([]*models.CheckIn, error) {
	return service.checkins.GetAllInRange(
		ctx,
		locationIDs,
		startDate,
		endDate,
	)
}

func (service CheckInService) GetByID(
	ctx context.Context,
	location *models.Location,
	id int64,
) (*models.CheckIn, error) {
	return service.checkins.GetByID(ctx, location, id)
}

func (service CheckInService) Create(
	ctx context.Context,
	location *models.Location,
	school *models.School,
) (*models.CheckIn, error) {
	checkIn, err := service.checkins.Create(ctx, location, school)
	if err != nil {
		return nil, err
	}

	location.Available--

	return checkIn, nil
}

func (service CheckInService) Delete(ctx context.Context, id int64) error {
	return service.checkins.Delete(ctx, id)
}
