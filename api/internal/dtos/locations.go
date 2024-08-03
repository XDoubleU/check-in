package dtos

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"github.com/xdoubleu/essentia/pkg/validate"

	"check-in/api/internal/models"
)

type CapacityMap = *orderedmap.OrderedMap[string, int64]
type SchoolsMap = *orderedmap.OrderedMap[string, int]

type CheckInsLocationEntryRaw struct {
	Capacities CapacityMap `json:"capacities" swaggertype:"object,number"`
	Schools    SchoolsMap  `json:"schools"    swaggertype:"object,number"`
} //	@name	CheckInsLocationEntryRaw

type PaginatedLocationsDto struct {
	PaginatedResultDto[models.Location]
} //	@name	PaginatedLocationsDto

type CreateLocationDto struct {
	Name             string            `json:"name"`
	Capacity         int64             `json:"capacity"`
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	TimeZone         string            `json:"timeZone"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	CreateLocationDto

type UpdateLocationDto struct {
	Name             *string           `json:"name"`
	Capacity         *int64            `json:"capacity"`
	Username         *string           `json:"username"`
	Password         *string           `json:"password"`
	TimeZone         *string           `json:"timeZone"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	UpdateLocationDto

func (dto *CreateLocationDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Name, validate.IsNotEmpty, "name")
	validate.Check(v, dto.Capacity, validate.IsGreaterThanFunc(int64(0)), "capacity")
	validate.Check(v, dto.Username, validate.IsNotEmpty, "username")
	validate.Check(v, dto.Password, validate.IsNotEmpty, "password")
	validate.Check(v, dto.TimeZone, validate.IsNotEmpty, "timeZone")
	validate.Check(v, dto.TimeZone, validate.IsValidTimeZone, "timeZone")

	dto.ValidationErrors = v.Errors

	return v
}

func (dto *UpdateLocationDto) Validate() *validate.Validator {
	v := validate.New()

	if dto.Name != nil {
		validate.Check(v, *dto.Name, validate.IsNotEmpty, "name")
	}

	if dto.Capacity != nil {
		validate.Check(
			v,
			*dto.Capacity,
			validate.IsGreaterThanFunc(int64(0)),
			"capacity",
		)
	}

	if dto.Username != nil {
		validate.Check(v, *dto.Username, validate.IsNotEmpty, "username")
	}

	if dto.Password != nil {
		validate.Check(v, *dto.Password, validate.IsNotEmpty, "password")
	}

	if dto.TimeZone != nil {
		validate.Check(v, *dto.TimeZone, validate.IsNotEmpty, "timeZone")
		validate.Check(v, *dto.TimeZone, validate.IsValidTimeZone, "timeZone")
	}

	dto.ValidationErrors = v.Errors

	return v
}
