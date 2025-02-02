package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"
	"github.com/jackc/pgx/v5/pgtype"
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

func (dto *CreateCheckInDto) Validate() (bool, map[string]string) {
	v := validate.New()

	validate.Check(v, "schoolId", dto.SchoolID, validate.IsGreaterThan(int64(0)))

	return v.Valid(), v.Errors()
}
