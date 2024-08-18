package dtos

import "github.com/XDoubleU/essentia/pkg/validate"

type SignInDto struct {
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	RememberMe       bool              `json:"rememberMe"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	SignInDto

func (dto *SignInDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.Username, validate.IsNotEmpty, "username")
	validate.Check(v, dto.Password, validate.IsNotEmpty, "password")

	dto.ValidationErrors = v.Errors

	return v
}
