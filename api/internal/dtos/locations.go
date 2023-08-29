package dtos

import (
	"strconv"
	"time"

	orderedmap "github.com/wk8/go-ordered-map/v2"

	"check-in/api/internal/models"
	"check-in/api/internal/validator"
)

type CapacityMap = *orderedmap.OrderedMap[string, int64]
type SchoolsMap = *orderedmap.OrderedMap[string, int]

type CheckInsLocationEntryRaw struct {
	Capacities CapacityMap `json:"capacities" swaggertype:"object,number"`
	Schools    SchoolsMap  `json:"schools"  swaggertype:"object,number"`
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

func ValidateCreateLocationDto(
	v *validator.Validator,
	createLocationDto CreateLocationDto,
) {
	v.Check(createLocationDto.Name != "", "name", "must be provided")
	v.Check(createLocationDto.Capacity > 0, "capacity", "must be greater than zero")
	v.Check(createLocationDto.Username != "", "username", "must be provided")
	v.Check(createLocationDto.Password != "", "password", "must be provided")

	_, err := time.LoadLocation(createLocationDto.TimeZone)
	v.Check(
		createLocationDto.TimeZone != "" && err == nil,
		"timeZone",
		"must be provided and must be a valid IANA value",
	)
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

	if updateLocationDto.TimeZone != nil {
		_, err := time.LoadLocation(*updateLocationDto.TimeZone)
		v.Check(
			err == nil,
			"timeZone",
			"must be provided and must be a valid IANA value",
		)
	}
}
