package dtos

import (
	"check-in/api/internal/models"
	"check-in/api/internal/validator"
)

type PaginatedSchoolsDto struct {
	PaginatedResultDto[models.School]
} //	@name	PaginatedSchoolsDto

type SchoolDto struct {
	Name string `json:"name"`
} //	@name	SchoolDto

func ValidateSchoolDto(v *validator.Validator, schoolDto SchoolDto) {
	v.Check(schoolDto.Name != "", "name", "must be provided")
}
