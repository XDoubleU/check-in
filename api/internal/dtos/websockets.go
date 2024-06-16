package dtos

import (
	"check-in/api/internal/models"

	"github.com/XDoubleU/essentia/pkg/validator"
)

type SubscribeMessageDto struct {
	Subject        models.WebSocketSubject `json:"subject"`
	NormalizedName string                  `json:"normalizedName"`
} //	@name	SubscribeMessageDto

func (dto SubscribeMessageDto) GetSubject() string {
	return string(dto.Subject)
}

func (dto SubscribeMessageDto) Validate() *validator.Validator {
	v := validator.New()

	if dto.Subject == models.SingleLocation {
		validator.Check(v, dto.NormalizedName, validator.IsNotEmpty, "normalizedName")
	}

	return v
}
