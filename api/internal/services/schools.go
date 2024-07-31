package services

import (
	"context"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
)

type SchoolService struct {
	schools repositories.SchoolRepository
}

func (service SchoolService) GetTotalCount(ctx context.Context) (*int64, error) {
	return service.schools.GetTotalCount(ctx)
}

// todo: refactor
func (service SchoolService) GetSchoolMaps(
	schools []*models.School,
) (map[int64]string, *orderedmap.OrderedMap[string, int]) {
	schoolsIDNameMap := make(map[int64]string)
	schoolsMap := orderedmap.New[string, int]()
	for _, school := range schools {
		schoolsIDNameMap[school.ID] = school.Name
		schoolsMap.Set(school.Name, 0)
	}

	return schoolsIDNameMap, schoolsMap
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
		return nil, errors.ErrFailedValidation
	}

	return service.schools.Create(ctx, schoolDto.Name)
}

func (service SchoolService) Update(
	ctx context.Context,
	id int64,
	schoolDto *dtos.SchoolDto,
) (*models.School, error) {
	if v := schoolDto.Validate(); !v.Valid() {
		return nil, errors.ErrFailedValidation
	}

	school, err := service.GetByIDWithoutReadOnly(ctx, id)
	if err != nil {
		return nil, errors.ErrResourceNotFound
	}

	err = service.schools.Update(ctx, school, schoolDto)
	if err != nil {
		return nil, err
	}

	return school, nil
}

func (service SchoolService) Delete(ctx context.Context, id int64) (*models.School, error) {
	school, err := service.GetByIDWithoutReadOnly(ctx, id)
	if err != nil {
		return nil, errors.ErrResourceNotFound
	}

	err = service.schools.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return school, nil
}
