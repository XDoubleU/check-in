package services

import (
	"context"
	"errors"

	"github.com/xdoubleu/essentia/pkg/database"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type SchoolService struct {
	schools         repositories.SchoolRepository
	schoolIDNameMap map[int64]string
}

func (service SchoolService) GetTotalCount(ctx context.Context) (*int64, error) {
	return service.schools.GetTotalCount(ctx)
}

func (service SchoolService) SchoolIDNameMap(
	ctx context.Context,
) (map[int64]string, error) {
	if len(service.schoolIDNameMap) != 0 {
		return service.schoolIDNameMap, nil
	}

	schools, err := service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, school := range schools {
		service.schoolIDNameMap[school.ID] = school.Name
	}

	return service.schoolIDNameMap, nil
}

func (service SchoolService) GetAll(ctx context.Context) ([]*models.School, error) {
	return service.schools.GetAll(ctx)
}

func (service SchoolService) GetAllSortedByLocation(
	ctx context.Context,
	locationID string,
) ([]*models.School, error) {
	return service.schools.GetAllSortedByLocation(ctx, locationID)
}

func (service SchoolService) GetAllPaginated(
	ctx context.Context,
	_ *models.User,
	limit int64,
	offset int64,
) ([]*models.School, error) {
	return service.schools.GetAllPaginated(ctx, limit, offset)
}

func (service SchoolService) GetByID(
	ctx context.Context,
	id int64,
) (*models.School, error) {
	return service.schools.GetByID(ctx, id)
}

func (service SchoolService) GetByIDWithoutReadOnly(
	ctx context.Context,
	id int64,
) (*models.School, error) {
	return service.schools.GetByIDWithoutReadOnly(ctx, id)
}

func (service SchoolService) Create(
	ctx context.Context,
	schoolDto *dtos.SchoolDto,
) (*models.School, error) {
	if v := schoolDto.Validate(); !v.Valid() {
		return nil, errortools.ErrFailedValidation
	}

	school, err := service.schools.Create(ctx, schoolDto.Name)
	if err != nil {
		if errors.Is(err, database.ErrResourceConflict) {
			return nil, errortools.NewConflictError("school", schoolDto.Name, "name")
		}
		return nil, err
	}

	service.schoolIDNameMap[school.ID] = school.Name

	return school, nil
}

func (service SchoolService) Update(
	ctx context.Context,
	id int64,
	schoolDto *dtos.SchoolDto,
) (*models.School, error) {
	if v := schoolDto.Validate(); !v.Valid() {
		return nil, errortools.ErrFailedValidation
	}

	school, err := service.GetByIDWithoutReadOnly(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("school", id, "id")
		}
		return nil, err
	}

	school, err = service.schools.Update(ctx, *school, schoolDto)
	if err != nil {
		if errors.Is(err, database.ErrResourceConflict) {
			return nil, errortools.NewConflictError("school", schoolDto.Name, "name")
		}
		return nil, err
	}

	service.schoolIDNameMap[school.ID] = school.Name

	return school, nil
}

func (service SchoolService) Delete(
	ctx context.Context,
	id int64,
) (*models.School, error) {
	school, err := service.GetByIDWithoutReadOnly(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrResourceNotFound) {
			return nil, errortools.NewNotFoundError("school", id, "id")
		}
		return nil, err
	}

	err = service.schools.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	delete(service.schoolIDNameMap, school.ID)

	return school, nil
}
