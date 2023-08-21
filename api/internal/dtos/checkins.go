package dtos

import (
	"check-in/api/internal/validator"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCheckInDto struct {
	SchoolID int64 `json:"schoolId"`
} //	@name	CreateCheckInDto

type CheckInDto struct {
	ID         int64              `json:"id"`
	LocationID string             `json:"locationId"`
	SchoolName   string              `json:"schoolName"`
	Capacity   int64              `json:"capacity"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"  swaggertype:"string"`
} // @name CheckInDto

func ValidateCreateCheckInDto(v *validator.Validator, CreatecheckInDto CreateCheckInDto) {
	v.Check(CreatecheckInDto.SchoolID > 0, "schoolId", "must be greater than zero")
}
