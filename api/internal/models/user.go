package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Roles string //	@name	Roles

const (
	DefaultRole Roles = "default"
	ManagerRole Roles = "manager"
	AdminRole   Roles = "admin"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash []byte `json:"-"`
	Role         Roles  `json:"role"`
} //	@name	User

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword(
		[]byte(password),
		12, //nolint:gomnd //no magic number
	)
}

func (user *User) CompareHashAndPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
