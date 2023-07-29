package dtos

import (
	"strconv"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"

	"check-in/api/internal/constants"
	"check-in/api/internal/models"
	"check-in/api/internal/validator"
)

type CheckInsLocationEntryRaw struct {
	Capacity int64                               `json:"capacity"`
	Schools  *orderedmap.OrderedMap[string, int] `json:"schools"`
} //	@name	CheckInsLocationEntryRaw

func ConvertCheckInsLocationEntryRawMapToCSV(
	timezone *time.Location,
	timeFormat string,
	entries *orderedmap.OrderedMap[int64, *CheckInsLocationEntryRaw],
) [][]string {
	var output [][]string

	var headers []string
	headers = append(headers, "datetime")
	headers = append(headers, "capacity")

	singleEntry := entries.Oldest().Value
	for school := singleEntry.Schools.Oldest(); school != nil; school = school.Next() {
		headers = append(headers, school.Key)
	}
	output = append(output, headers)

	for pair := entries.Oldest(); pair != nil; pair = pair.Next() {
		var entry []string

		entry = append(
			entry,
			time.Unix(pair.Key/constants.SecToMilliSec, 0).
				In(timezone).
				Format(timeFormat),
		)
		entry = append(entry, strconv.FormatInt(pair.Value.Capacity, 10))

		for school := pair.Value.Schools.Oldest(); school != nil; school = school.Next() {
			entry = append(entry, strconv.Itoa(school.Value))
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
