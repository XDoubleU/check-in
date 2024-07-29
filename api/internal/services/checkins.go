package services

import (
	"context"
	"time"

	"github.com/xdoubleu/essentia/pkg/tools"

	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type CheckInService struct {
	checkins  repositories.CheckInRepository
	websocket *WebSocketService
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

	service.websocket.NewLocationState(*location)

	return checkIn, nil
}

func (service CheckInService) Delete(ctx context.Context, id int64) error {
	return service.checkins.Delete(ctx, id)
}
