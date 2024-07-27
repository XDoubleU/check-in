package dtos

import (
	"strconv"

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

// todo refactor
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

func (dto CreateLocationDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Name, validate.IsNotEmpty, "name")
	validate.Check(v, dto.Capacity, validate.IsGreaterThanFunc(int64(0)), "capacity")
	validate.Check(v, dto.Username, validate.IsNotEmpty, "username")
	validate.Check(v, dto.Password, validate.IsNotEmpty, "password")
	validate.Check(v, dto.TimeZone, validate.IsNotEmpty, "timeZone")
	validate.Check(v, dto.TimeZone, validate.IsValidTimeZone, "timeZone")

	return v
}

func (dto UpdateLocationDto) Validate() *validate.Validator {
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

	return v
}
