package dtos

import "github.com/XDoubleU/essentia/pkg/validate"

type SignInDto struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
} //	@name	SignInDto

func (dto *SignInDto) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.Check(v, "username", dto.Username, validate.IsNotEmpty)
	validate.Check(v, "password", dto.Password, validate.IsNotEmpty)

	return v.Valid(), v.Errors()
}
