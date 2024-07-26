package dtos

import (
	"github.com/xdoubleu/essentia/pkg/validate"

	"check-in/api/internal/models"
)

type SubscribeMessageDto struct {
	TopicName      models.WebSocketTopic `json:"topic"`
	NormalizedName string                `json:"normalizedName"`
} //	@name	SubscribeMessageDto

func (dto SubscribeMessageDto) Topic() string {
	return string(dto.TopicName)
}

func (dto SubscribeMessageDto) Validate() *validate.Validator {
	v := validate.New()

	if dto.TopicName == models.SingleLocation {
		validate.Check(v, dto.NormalizedName, validate.IsNotEmpty, "normalizedName")
	}

	return v
}
