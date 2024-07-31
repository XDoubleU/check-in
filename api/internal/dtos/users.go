package dtos

import (
	"github.com/xdoubleu/essentia/pkg/validate"

	"check-in/api/internal/models"
)

type PaginatedUsersDto struct {
	PaginatedResultDto[models.User]
} //	@name	PaginatedUsersDto

type CreateUserDto struct {
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	CreateUserDto

type UpdateUserDto struct {
	Username         *string           `json:"username"`
	Password         *string           `json:"password"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	UpdateUserDto

func (dto *CreateUserDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Username, validate.IsNotEmpty, "username")
	validate.Check(v, dto.Password, validate.IsNotEmpty, "password")

	dto.ValidationErrors = v.Errors

	return v
}

func (dto *UpdateUserDto) Validate() *validate.Validator {
	v := validate.New()

	if dto.Username != nil {
		validate.Check(v, *dto.Username, validate.IsNotEmpty, "username")
	}

	if dto.Password != nil {
		validate.Check(v, *dto.Password, validate.IsNotEmpty, "password")
	}

	dto.ValidationErrors = v.Errors

	return v
}
