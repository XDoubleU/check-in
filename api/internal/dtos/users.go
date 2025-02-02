package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"

	"check-in/api/internal/models"
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

func (dto *CreateUserDto) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.Check(v, "username", dto.Username, validate.IsNotEmpty)
	validate.Check(v, "password", dto.Password, validate.IsNotEmpty)

	return v.Valid(), v.Errors()
}

func (dto *UpdateUserDto) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.CheckOptional(v, "username", dto.Username, validate.IsNotEmpty)
	validate.CheckOptional(v, "password", dto.Password, validate.IsNotEmpty)

	return v.Valid(), v.Errors()
}
