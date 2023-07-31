package dtos

import (
	"check-in/api/internal/models"
	"check-in/api/internal/validator"
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

func ValidateCreateUserDto(v *validator.Validator, createUserDto CreateUserDto) {
	v.Check(createUserDto.Username != "", "username", "must be provided")
	v.Check(createUserDto.Password != "", "password", "must be provided")
}

func ValidateUpdateUserDto(v *validator.Validator, updateUserDto UpdateUserDto) {
	if updateUserDto.Username != nil {
		v.Check(*updateUserDto.Username != "", "username", "must be provided")
	}

	if updateUserDto.Password != nil {
		v.Check(*updateUserDto.Password != "", "password", "must be provided")
	}
}
