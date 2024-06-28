package dtos

import (
	"github.com/XDoubleU/essentia/pkg/validate"

	"check-in/api/internal/models"
)

type SubscribeMessageDto struct {
	Subject        models.WebSocketSubject `json:"subject"`
	NormalizedName string                  `json:"normalizedName"`
} //	@name	SubscribeMessageDto

func (dto SubscribeMessageDto) GetSubject() string {
	return string(dto.Subject)
}

func (dto SubscribeMessageDto) Validate() *validate.Validator {
	v := validate.New()

	if dto.Subject == models.SingleLocation {
		validate.Check(v, dto.NormalizedName, validate.IsNotEmpty, "normalizedName")
	}

	return v
}
