package dtos

import (
	"strconv"

	"github.com/XDoubleU/essentia/pkg/validator"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"check-in/api/internal/models"
)

type CapacityMap = *orderedmap.OrderedMap[string, int64]
type SchoolsMap = *orderedmap.OrderedMap[string, int]

type CheckInsLocationEntryRaw struct {
	Capacities CapacityMap `json:"capacities" swaggertype:"object,number"`
	Schools    SchoolsMap  `json:"schools"    swaggertype:"object,number"`
} //	@name	CheckInsLocationEntryRaw

func ConvertCheckInsLocationEntryRawMapToCSV(
	entries *orderedmap.OrderedMap[string, *CheckInsLocationEntryRaw],
) [][]string {
	var output [][]string

	var headers []string
	headers = append(headers, "datetime")
	headers = append(headers, "capacity")

	if entries.Len() == 0 {
		output = append(output, headers)
		return output
	}

	singleEntry := entries.Oldest().Value
	for school := singleEntry.Schools.Oldest(); school != nil; school = school.Next() {
		headers = append(headers, school.Key)
	}
	output = append(output, headers)

	for pair := entries.Oldest(); pair != nil; pair = pair.Next() {
		var entry []string

		var totalCapacity int64
		capacities := pair.Value.Capacities
		for capacity := capacities.Oldest(); capacity != nil; capacity = capacity.Next() {
			totalCapacity += capacity.Value
		}

		entry = append(
			entry,
			pair.Key,
		)
		entry = append(entry, strconv.FormatInt(totalCapacity, 10))

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
	TimeZone string `json:"timeZone"`
} //	@name	CreateLocationDto

type UpdateLocationDto struct {
	Name     *string `json:"name"`
	Capacity *int64  `json:"capacity"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	TimeZone *string `json:"timeZone"`
} //	@name	UpdateLocationDto

func (dto CreateLocationDto) Validate() *validator.Validator {
	v := validator.New()

	validator.Check(v, dto.Name, validator.IsNotEmpty, "name")
	validator.Check(v, dto.Capacity, validator.IsGreaterThanFunc(int64(0)), "capacity")
	validator.Check(v, dto.Username, validator.IsNotEmpty, "username")
	validator.Check(v, dto.Password, validator.IsNotEmpty, "password")
	validator.Check(v, dto.TimeZone, validator.IsNotEmpty, "timeZone")
	validator.Check(v, dto.TimeZone, validator.IsValidTimeZone, "timeZone")

	return v
}

func (dto UpdateLocationDto) Validate() *validator.Validator {
	v := validator.New()

	if dto.Name != nil {
		validator.Check(v, *dto.Name, validator.IsNotEmpty, "name")
	}

	if dto.Capacity != nil {
		validator.Check(v, *dto.Capacity, validator.IsGreaterThanFunc(int64(0)), "capacity")
	}

	if dto.Username != nil {
		validator.Check(v, *dto.Username, validator.IsNotEmpty, "username")
	}

	if dto.Password != nil {
		validator.Check(v, *dto.Password, validator.IsNotEmpty, "password")
	}

	if dto.TimeZone != nil {
		validator.Check(v, *dto.TimeZone, validator.IsNotEmpty, "timeZone")
		validator.Check(v, *dto.TimeZone, validator.IsValidTimeZone, "timeZone")
	}

	return v
}
