package dtos

import "check-in/api/internal/validator"

type SignInDto struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
} //	@name	SignInDto

func ValidateSignInDto(v *validator.Validator, signInDto SignInDto) {
	v.Check(signInDto.Username != "", "username", "must be provided")
	v.Check(signInDto.Password != "", "password", "must be provided")
}
