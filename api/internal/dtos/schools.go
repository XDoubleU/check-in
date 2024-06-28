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

func (dto SchoolDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Name, validate.IsNotEmpty, "name")

	return v
}
