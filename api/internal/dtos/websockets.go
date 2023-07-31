package dtos

import (
	"check-in/api/internal/models"
	"check-in/api/internal/validator"
)

type SubscribeMessageDto struct {
	Subject        models.WebSocketSubject `json:"subject"`
	NormalizedName string                  `json:"normalizedName"`
} //	@name	SubscribeMessageDto

func ValidateSubscribeMessageDto(
	v *validator.Validator,
	subscribeMessageDto SubscribeMessageDto,
) {
	if subscribeMessageDto.Subject == models.SingleLocation {
		v.Check(
			subscribeMessageDto.NormalizedName != "",
			"normalizedName",
			"must be provided",
		)
	}
}
