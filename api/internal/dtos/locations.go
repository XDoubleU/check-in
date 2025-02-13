package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"

	"check-in/api/internal/models"
)

type CheckInsGraphDto struct {
	Dates                 []string         `json:"dates"`
	CapacitiesPerLocation map[string][]int `json:"capacitiesPerLocation"`
	ValuesPerSchool       map[string][]int `json:"valuesPerSchool"`
} //	@name	CheckInsGraphDto

type PaginatedLocationsDto struct {
	PaginatedResultDto[models.Location]
} //	@name	PaginatedLocationsDto

type CreateLocationDto struct {
	Name     string `json:"name"`
	Capacity int64  `json:"capacity"`
	Username string `json:"username"`
	Password string `json:"password"`
	TimeZone string `json:"timeZone"`
} //	@name	CreateLocationDto

type UpdateLocationDto struct {
	Name     *string `json:"name"`
	Capacity *int64  `json:"capacity"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	TimeZone *string `json:"timeZone"`
} //	@name	UpdateLocationDto

func (dto *CreateLocationDto) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.Check(v, "name", dto.Name, validate.IsNotEmpty)
	validate.Check(v, "capacity", dto.Capacity, validate.IsGreaterThan(int64(0)))
	validate.Check(v, "username", dto.Username, validate.IsNotEmpty)
	validate.Check(v, "password", dto.Password, validate.IsNotEmpty)
	validate.Check(v, "timeZone", dto.TimeZone, validate.IsNotEmpty)
	validate.Check(v, "timeZone", dto.TimeZone, validate.IsValidTimeZone)

	return v.Valid(), v.Errors()
}

func (dto *UpdateLocationDto) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.CheckOptional(v, "name", dto.Name, validate.IsNotEmpty)
	validate.CheckOptional(
		v,
		"capacity",
		dto.Capacity,
		validate.IsGreaterThan(int64(0)),
	)
	validate.CheckOptional(v, "username", dto.Username, validate.IsNotEmpty)
	validate.CheckOptional(v, "password", dto.Password, validate.IsNotEmpty)
	validate.CheckOptional(v, "timeZone", dto.TimeZone, validate.IsNotEmpty)
	validate.CheckOptional(v, "timeZone", dto.TimeZone, validate.IsValidTimeZone)

	return v.Valid(), v.Errors()
}
