package models

import "time"

type CheckIn struct {
	ID         int64     `json:"id"`
	LocationID string    `json:"locationId"`
	SchoolID   int64     `json:"schoolId"`
	Capacity   int64     `json:"capacity"`
	CreatedAt  time.Time `json:"createdAt"`
} //	@name	CheckIn
