package dtos

import (
	"check-in/api/internal/models"

	"github.com/XDoubleU/essentia/pkg/validator"
)

type PaginatedSchoolsDto struct {
	PaginatedResultDto[models.School]
} //	@name	PaginatedSchoolsDto

type SchoolDto struct {
	Name string `json:"name"`
} //	@name	SchoolDto

func (dto SchoolDto) Validate() *validator.Validator {
	v := validator.New()

	validator.Check(v, dto.Name, validator.IsNotEmpty, "name")

	return v
}
