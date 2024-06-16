package dtos

import "github.com/XDoubleU/essentia/pkg/validator"

type SignInDto struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
} //	@name	SignInDto

func (dto SignInDto) Validate() *validator.Validator {
	v := validator.New()

	validator.Check(v, dto.Username, validator.IsNotEmpty, "username")
	validator.Check(v, dto.Password, validator.IsNotEmpty, "password")

	return v
}
