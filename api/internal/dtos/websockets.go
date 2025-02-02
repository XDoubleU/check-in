package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"
	"github.com/jackc/pgx/v5/pgtype"

	"check-in/api/internal/models"
)

type WebSocketSubject string //	@name	WebSocketSubject

const (
	AllLocations   WebSocketSubject = "all-locations"
	SingleLocation WebSocketSubject = "single-location"
	State          WebSocketSubject = "state"
)

type LocationStateDto struct {
	NormalizedName     string             `json:"normalizedName"`
	Available          int64              `json:"available"`
	Capacity           int64              `json:"capacity"`
	AvailableYesterday int64              `json:"availableYesterday"`
	CapacityYesterday  int64              `json:"capacityYesterday"`
	YesterdayFullAt    pgtype.Timestamptz `json:"yesterdayFullAt"    swaggertype:"string"`
} //	@name	LocationUpdateEvent

func NewLocationStateDto(location models.Location) LocationStateDto {
	return LocationStateDto{
		NormalizedName:     location.NormalizedName,
		Available:          location.Available,
		Capacity:           location.Capacity,
		YesterdayFullAt:    location.YesterdayFullAt,
		AvailableYesterday: location.AvailableYesterday,
		CapacityYesterday:  location.CapacityYesterday,
	}
}

type SubscribeMessageDto struct {
	Subject        WebSocketSubject `json:"subject"`
	NormalizedName string           `json:"normalizedName"`
} //	@name	SubscribeMessageDto

func (dto SubscribeMessageDto) Topic() string {
	if dto.Subject == AllLocations {
		return "*"
	}

	if dto.Subject == SingleLocation {
		return dto.NormalizedName
	}

	return string(dto.Subject)
}

func (dto SubscribeMessageDto) Validate() (bool, map[string]string) {
	v := validate.New()

	if dto.Subject == SingleLocation {
		validate.Check(v, "normalizedName", dto.NormalizedName, validate.IsNotEmpty)
	}

	return v.Valid(), v.Errors()
}
