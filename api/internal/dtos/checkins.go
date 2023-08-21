package dtos

import (
	"github.com/jackc/pgx/v5/pgtype"

	"check-in/api/internal/validator"
)

type CreateCheckInDto struct {
	SchoolID int64 `json:"schoolId"`
} //	@name	CreateCheckInDto

type CheckInDto struct {
	ID         int64              `json:"id"`
	LocationID string             `json:"locationId"`
	SchoolName string             `json:"schoolName"`
	Capacity   int64              `json:"capacity"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"  swaggertype:"string"`
} //	@name	CheckInDto

func ValidateCreateCheckInDto(
	v *validator.Validator,
	createCheckInDto CreateCheckInDto,
) {
	v.Check(createCheckInDto.SchoolID > 0, "schoolId", "must be greater than zero")
}
