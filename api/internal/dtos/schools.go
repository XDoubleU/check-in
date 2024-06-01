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

	v.Check(dto.Name != "", "name", "must be provided")

	return v
}
