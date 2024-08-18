package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCheckInDto struct {
	SchoolID         int64             `json:"schoolId"`
	ValidationErrors map[string]string `json:"-"`
} //	@name	CreateCheckInDto

type CheckInDto struct {
	ID         int64              `json:"id"`
	LocationID string             `json:"locationId"`
	SchoolName string             `json:"schoolName"`
	Capacity   int64              `json:"capacity"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"  swaggertype:"string"`
} //	@name	CheckInDto

func (dto *CreateCheckInDto) Validate() *validate.Validator {
	v := validate.New()

	validate.Check(v, dto.SchoolID, validate.IsGreaterThanFunc(int64(0)), "schoolId")

	dto.ValidationErrors = v.Errors

	return v
}
