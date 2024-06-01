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

	v.Check(dto.Username != "", "username", "must be provided")
	v.Check(dto.Password != "", "password", "must be provided")

	return v
}

func (dto UpdateUserDto) Validate() *validator.Validator {
	v := validator.New()

	if dto.Username != nil {
		v.Check(*dto.Username != "", "username", "must be provided")
	}

	if dto.Password != nil {
		v.Check(*dto.Password != "", "password", "must be provided")
	}

	return v
}
