package services

import (
	"context"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"

	"github.com/xdoubleu/essentia/pkg/errors"
)

type CheckInWriterService struct {
	checkins  repositories.CheckInWriteRepository
	locations LocationService
	schools   SchoolService
}

func (service CheckInWriterService) GetAllSchoolsSortedByLocation(
	ctx context.Context,
	userID string,
) ([]*models.School, error) {
	location, err := service.locations.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrResourceNotFound
	}

	return service.schools.GetAllSortedByLocation(ctx, location.ID)
}

func (service CheckInWriterService) Create(
	ctx context.Context,
	createCheckInDto *dtos.CreateCheckInDto,
	user *models.User,
) (*dtos.CheckInDto, error) {
	if v := createCheckInDto.Validate(); !v.Valid() {
		return nil, errors.ErrFailedValidation
	}

	location, err := service.locations.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	school, err := service.schools.GetByID(
		ctx,
		createCheckInDto.SchoolID,
	)
	if err != nil {
		return nil, errors.ErrResourceNotFound
	}

	if location.Available <= 0 {
		return nil, errors.ErrBadRequest
	}

	checkIn, err := service.checkins.Create(ctx, location, school)
	if err != nil {
		return nil, err
	}

	service.locations.NewCheckIn(*location)

	checkInDto := &dtos.CheckInDto{
		ID:         checkIn.ID,
		LocationID: checkIn.LocationID,
		SchoolName: school.Name,
		Capacity:   checkIn.Capacity,
		CreatedAt:  checkIn.CreatedAt,
	}

	return checkInDto, nil
}
