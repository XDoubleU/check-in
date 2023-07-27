package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CheckIn struct {
	ID         int64              `json:"id"`
	LocationID string             `json:"locationId"`
	SchoolID   int64              `json:"schoolId"`
	Capacity   int64              `json:"capacity"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"  swaggertype:"string"`
} //	@name	CheckIn
