package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type WebSocketTopic string //	@name	WebSocketSubject

const (
	AllLocations   WebSocketTopic = "all-locations"
	SingleLocation WebSocketTopic = "single-location"
)

type LocationUpdateEvent struct {
	NormalizedName     string             `json:"normalizedName"`
	Available          int64              `json:"available"`
	Capacity           int64              `json:"capacity"`
	AvailableYesterday int64              `json:"availableYesterday"`
	CapacityYesterday  int64              `json:"capacityYesterday"`
	YesterdayFullAt    pgtype.Timestamptz `json:"yesterdayFullAt"    swaggertype:"string"`
} //	@name	LocationUpdateEvent
