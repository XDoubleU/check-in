package services

import (
	"context"
	"net/http"

	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	errortools "github.com/xdoubleu/essentia/pkg/errors"

	"check-in/api/internal/dtos"
	"check-in/api/internal/models"
)

type WebSocketService struct {
	handler   *wstools.WebSocketHandler[dtos.SubscribeMessageDto]
	allTopic  *wstools.Topic
	topics    map[string]*wstools.Topic
	locations LocationService
}

func NewWebSocketService(
	allowedOrigin string,
	locationService LocationService,
) (*WebSocketService, error) {
	service := &WebSocketService{
		handler:   nil,
		allTopic:  nil,
		topics:    make(map[string]*wstools.Topic),
		locations: locationService,
	}

	handler := wstools.CreateWebSocketHandler[dtos.SubscribeMessageDto](
		1,   //nolint:mnd //no magic number
		100, //nolint:mnd //no magic number
		[]string{allowedOrigin},
	)
	service.handler = &handler

	locations, err := service.locations.GetAll(context.Background())
	if err != nil {
		return nil, err
	}

	service.allTopic, err = service.handler.AddTopic("*", func(topic *wstools.Topic) any { return service.getAllLocationStates() })
	if err != nil {
		return nil, err
	}

	for _, location := range locations {
		err = service.AddLocation(location)
		if err != nil {
			return nil, err
		}
	}

	return service, nil
}

func (service WebSocketService) Handler() http.HandlerFunc {
	return service.handler.Handler()
}

func (service WebSocketService) AddLocation(location *models.Location) error {
	topic, err := service.handler.AddTopic(location.NormalizedName, nil)
	if err != nil {
		return err
	}

	service.topics[location.ID] = topic
	return nil
}

func (service WebSocketService) UpdateLocation(location *models.Location) error {
	err := service.DeleteLocation(location)
	if err != nil {
		return err
	}

	return service.AddLocation(location)
}

func (service WebSocketService) DeleteLocation(location *models.Location) error {
	topic, ok := service.topics[location.ID]
	if !ok {
		return errortools.ErrResourceNotFound
	}

	err := service.handler.RemoveTopic(topic)
	if err != nil {
		return err
	}

	delete(service.topics, topic.Name)
	return nil
}

func (service WebSocketService) NewLocationState(location models.Location) {
	locationState := dtos.NewLocationStateDto(location)

	service.allTopic.EnqueueEvent(locationState)
	service.topics[location.ID].EnqueueEvent(locationState)
}

func (service WebSocketService) getAllLocationStates() []dtos.LocationStateDto {
	locations, err := service.locations.GetAll(context.Background())
	if err != nil {
		//todo no panic pls
		panic(err)
	}

	var result []dtos.LocationStateDto
	for _, location := range locations {
		result = append(result, dtos.NewLocationStateDto(*location))
	}

	return result
}
