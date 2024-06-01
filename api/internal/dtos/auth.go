package dtos

import "github.com/XDoubleU/essentia/pkg/validator"

type SignInDto struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
} //	@name	SignInDto

func (dto SignInDto) Validate() *validator.Validator {
	v := validator.New()

	v.Check(dto.Username != "", "username", "must be provided")
	v.Check(dto.Password != "", "password", "must be provided")

	return v
}
