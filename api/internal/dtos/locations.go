package dtos

import (
	"encoding/json"

	"github.com/XDoubleU/essentia/pkg/validate"
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"check-in/api/internal/models"
)

type CapacityMap = *orderedmap.OrderedMap[string, int64]
type SchoolsMap = *orderedmap.OrderedMap[string, int]

type CheckInsLocationEntryRaw struct {
	Capacities CapacityMap `json:"capacities" swaggertype:"object,number"`
	Schools    SchoolsMap  `json:"schools"    swaggertype:"object,number"`
} //	@name	CheckInsLocationEntryRaw

func NewCheckInsLocationEntryRaw(
	capacities CapacityMap,
	schools SchoolsMap,
) CheckInsLocationEntryRaw {
	// Used to deep copy capacities
	var capacitiesCopy *orderedmap.OrderedMap[string, int64]
	data, _ := json.Marshal(capacities)
	_ = json.Unmarshal(data, &capacitiesCopy)

	// Used to deep copy schools
	var schoolsCopy *orderedmap.OrderedMap[string, int]
	data, _ = json.Marshal(schools)
	_ = json.Unmarshal(data, &schoolsCopy)

	return CheckInsLocationEntryRaw{
		Capacities: capacitiesCopy,
		Schools:    schoolsCopy,
	}
}

func (entry CheckInsLocationEntryRaw) Copy() CheckInsLocationEntryRaw {
	return CheckInsLocationEntryRaw{
		Capacities: entry.Capacities,
		Schools:    entry.Schools,
	}
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
