package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"

	"check-in/api/internal/models"
)

type PaginatedSchoolsDto struct {
	PaginatedResultDto[models.School]
} //	@name	PaginatedSchoolsDto

type SchoolDto struct {
	Name string `json:"name"`
} //	@name	SchoolDto

func (dto *SchoolDto) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.Check(v, "name", dto.Name, validate.IsNotEmpty)

	return v.Valid(), v.Errors()
}
