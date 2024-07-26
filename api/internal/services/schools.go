package services

import (
	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
	"check-in/api/internal/repositories"
	"context"

	orderedmap "github.com/wk8/go-ordered-map/v2"
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
	name string,
) (*models.School, error) {
	return service.schools.Create(ctx, name)
}

func (service SchoolService) Update(
	ctx context.Context,
	school *models.School,
	schoolDto dtos.SchoolDto,
) error {
	return service.schools.Update(ctx, school, schoolDto)
}

func (service SchoolService) Delete(ctx context.Context, id int64) error {
	return service.schools.Delete(ctx, id)
}
