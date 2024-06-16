package dtos

import (
	"check-in/api/internal/models"

	"github.com/XDoubleU/essentia/pkg/validator"
)

type PaginatedUsersDto struct {
	PaginatedResultDto[models.User]
} //	@name	PaginatedUsersDto

type CreateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
} //	@name	CreateUserDto

type UpdateUserDto struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
} //	@name	UpdateUserDto

func (dto CreateUserDto) Validate() *validator.Validator {
	v := validator.New()

	validator.Check(v, dto.Username, validator.IsNotEmpty, "username")
	validator.Check(v, dto.Password, validator.IsNotEmpty, "password")

	return v
}

func (dto UpdateUserDto) Validate() *validator.Validator {
	v := validator.New()

	if dto.Username != nil {
		validator.Check(v, *dto.Username, validator.IsNotEmpty, "username")
	}

	if dto.Password != nil {
		validator.Check(v, *dto.Password, validator.IsNotEmpty, "password")
	}

	return v
}
