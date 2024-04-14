package models

import (
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
)

type WebSocketSubject string //	@name	WebSocketSubject

const (
	AllLocations   WebSocketSubject = "all-locations"
	SingleLocation WebSocketSubject = "single-location"
)

type Subscriber struct {
	Subject        WebSocketSubject
	NormalizedName string
	Buffer         map[string]LocationUpdateEvent
	BufferMu       *sync.Mutex
}

type LocationUpdateEvent struct {
	NormalizedName     string             `json:"normalizedName"`
	Available          int64              `json:"available"`
	Capacity           int64              `json:"capacity"`
	AvailableYesterday int64              `json:"availableYesterday"`
	CapacityYesterday  int64              `json:"capacityYesterday"`
	YesterdayFullAt    pgtype.Timestamptz `json:"yesterdayFullAt"    swaggertype:"string"`
} //	@name	LocationUpdateEvent
