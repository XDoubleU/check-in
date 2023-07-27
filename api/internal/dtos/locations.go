package dtos

import (
	"check-in/api/internal/models"
	"check-in/api/internal/validator"
)

type CheckInsLocationEntryRaw struct {
	Capacity int64          `json:"capacity"`
	Schools  map[string]int `json:"schools"`
} //	@name	CheckInsLocationEntryRaw

type CheckInsLocationEntryCsv struct {
	Datetime int64          `json:"datetime"`
	Capacity int64          `json:"capacity"`
	Schools  map[string]int `json:"schools"`
} //	@name	CheckInsLocationEntryCsv

func ConvertCheckInsLocationEntryRawMapToCsv(
	entries map[int64]*CheckInsLocationEntryRaw,
) []*CheckInsLocationEntryCsv {
	var output []*CheckInsLocationEntryCsv
	for k, v := range entries {
		entry := &CheckInsLocationEntryCsv{
			Datetime: k,
			Capacity: v.Capacity,
			Schools:  v.Schools,
		}

		output = append(output, entry)
	}

	return output
}

type PaginatedLocationsDto struct {
	PaginatedResultDto[models.Location]
} //	@name	PaginatedLocationsDto

type CreateLocationDto struct {
	Name     string `json:"name"`
	Capacity int64  `json:"capacity"`
	Username string `json:"username"`
	Password string `json:"password"`
} //	@name	CreateLocationDto

type UpdateLocationDto struct {
	Name     *string `json:"name"`
	Capacity *int64  `json:"capacity"`
	Username *string `json:"username"`
	Password *string `json:"password"`
} //	@name	UpdateLocationDto

func ValidateCreateLocationDto(
	v *validator.Validator,
	createLocationDto CreateLocationDto,
) {
	v.Check(createLocationDto.Name != "", "name", "must be provided")
	v.Check(createLocationDto.Capacity > 0, "capacity", "must be greater than zero")
	v.Check(createLocationDto.Username != "", "username", "must be provided")
	v.Check(createLocationDto.Password != "", "password", "must be provided")
}

func ValidateUpdateLocationDto(
	v *validator.Validator,
	updateLocationDto UpdateLocationDto,
) {
	if updateLocationDto.Name != nil {
		v.Check(*updateLocationDto.Name != "", "name", "must be provided")
	}

	if updateLocationDto.Capacity != nil {
		v.Check(
			*updateLocationDto.Capacity > 0,
			"capacity",
			"must be greater than zero",
		)
	}

	if updateLocationDto.Username != nil {
		v.Check(*updateLocationDto.Username != "", "username", "must be provided")
	}

	if updateLocationDto.Password != nil {
		v.Check(*updateLocationDto.Password != "", "password", "must be provided")
	}
}
