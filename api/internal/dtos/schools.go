package dtos

import (
	"github.com/xdoubleu/essentia/pkg/validate"

	"check-in/api/internal/models"
)

type PaginatedSchoolsDto struct {
	PaginatedResultDto[models.School]
} //	@name	PaginatedSchoolsDto

type SchoolDto struct {
	Name             string            `json:"name"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	SchoolDto

func (dto *SchoolDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Name, validate.IsNotEmpty, "name")

	dto.ValidationErrors = v.Errors

	return v
}
