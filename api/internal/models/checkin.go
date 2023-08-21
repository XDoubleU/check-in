package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CheckIn struct {
	ID         int64
	LocationID string
	SchoolID   int64
	Capacity   int64
	CreatedAt  pgtype.Timestamptz
}
